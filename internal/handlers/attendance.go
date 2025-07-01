package handlers

import (
	"database/sql"
	"fmt"
	"labor-management-system/database"
	"labor-management-system/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ClockInRequest struct {
	EmployeeID int `json:"employee_id" binding:"required"`
}

type ClockOutRequest struct {
	EmployeeID int `json:"employee_id" binding:"required"`
}

func GetAttendanceLogs(c *gin.Context) {
	// Get query parameters for filtering
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	employeeID := c.Query("employee_id")

	query := `
		SELECT a.id, a.employee_id, a.work_date, a.clock_in, a.clock_out, 
		       a.break_start, a.break_end, a.total_hours, a.overtime_hours, 
		       a.status, a.notes, a.created_at, a.updated_at,
		       e.name as employee_name, e.employee_number
		FROM attendance_logs a
		JOIN employees e ON a.employee_id = e.id
		WHERE 1=1
	`

	args := []interface{}{}

	if startDate != "" {
		query += " AND a.work_date >= ?"
		args = append(args, startDate)
	}

	if endDate != "" {
		query += " AND a.work_date <= ?"
		args = append(args, endDate)
	}

	if employeeID != "" {
		query += " AND a.employee_id = ?"
		args = append(args, employeeID)
	}

	query += " ORDER BY a.work_date DESC, e.name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var attendances []map[string]interface{}
	for rows.Next() {
		var attendance models.AttendanceLog
		var employeeName, employeeNumber string

		err := rows.Scan(
			&attendance.ID, &attendance.EmployeeID, &attendance.WorkDate,
			&attendance.ClockIn, &attendance.ClockOut, &attendance.BreakStart,
			&attendance.BreakEnd, &attendance.TotalHours, &attendance.OvertimeHours,
			&attendance.Status, &attendance.Notes, &attendance.CreatedAt,
			&attendance.UpdatedAt, &employeeName, &employeeNumber,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan attendance"})
			return
		}

		attendanceData := map[string]interface{}{
			"attendance":      attendance,
			"employee_name":   employeeName,
			"employee_number": employeeNumber,
		}
		attendances = append(attendances, attendanceData)
	}

	c.JSON(http.StatusOK, gin.H{"attendances": attendances})
}

func GetEmployeeAttendance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	// Get query parameters for date range
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := `
		SELECT a.id, a.employee_id, a.work_date, a.clock_in, a.clock_out, 
		       a.break_start, a.break_end, a.total_hours, a.overtime_hours, 
		       a.status, a.notes, a.created_at, a.updated_at
		FROM attendance_logs a
		WHERE a.employee_id = ?
	`

	args := []interface{}{id}

	if startDate != "" {
		query += " AND a.work_date >= ?"
		args = append(args, startDate)
	}

	if endDate != "" {
		query += " AND a.work_date <= ?"
		args = append(args, endDate)
	}

	query += " ORDER BY a.work_date DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var attendances []models.AttendanceLog
	for rows.Next() {
		var attendance models.AttendanceLog

		err := rows.Scan(
			&attendance.ID, &attendance.EmployeeID, &attendance.WorkDate,
			&attendance.ClockIn, &attendance.ClockOut, &attendance.BreakStart,
			&attendance.BreakEnd, &attendance.TotalHours, &attendance.OvertimeHours,
			&attendance.Status, &attendance.Notes, &attendance.CreatedAt,
			&attendance.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan attendance"})
			return
		}

		attendances = append(attendances, attendance)
	}

	c.JSON(http.StatusOK, gin.H{"attendances": attendances})
}

func ClockIn(c *gin.Context) {
	var req ClockInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	today := time.Now().Format("2006-01-02")
	now := time.Now().Format("15:04:05")

	// Check if employee already clocked in today
	var existingID int
	err := database.DB.QueryRow(
		"SELECT id FROM attendance_logs WHERE employee_id = ? AND work_date = ?",
		req.EmployeeID, today,
	).Scan(&existingID)

	if err == nil {
		// Already has a record for today, update clock_in if not set
		var existingClockIn sql.NullString
		err = database.DB.QueryRow(
			"SELECT clock_in FROM attendance_logs WHERE id = ?",
			existingID,
		).Scan(&existingClockIn)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if existingClockIn.Valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Already clocked in today"})
			return
		}

		// Update clock_in time
		_, err = database.DB.Exec(
			"UPDATE attendance_logs SET clock_in = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
			now, existingID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update clock in"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Clocked in successfully", "time": now})
		return
	}

	if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Create new attendance record
	_, err = database.DB.Exec(`
		INSERT INTO attendance_logs (employee_id, work_date, clock_in, status)
		VALUES (?, ?, ?, 'present')
	`, req.EmployeeID, today, now)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clock in"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Clocked in successfully", "time": now})
}

func ClockOut(c *gin.Context) {
	var req ClockOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	today := time.Now().Format("2006-01-02")
	now := time.Now().Format("15:04:05")

	// Find today's attendance record
	var attendanceID int
	var clockInStr sql.NullString
	var clockOutStr sql.NullString

	err := database.DB.QueryRow(
		"SELECT id, clock_in, clock_out FROM attendance_logs WHERE employee_id = ? AND work_date = ?",
		req.EmployeeID, today,
	).Scan(&attendanceID, &clockInStr, &clockOutStr)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No clock in record found for today"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if !clockInStr.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Must clock in first"})
		return
	}

	if clockOutStr.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already clocked out today"})
		return
	}

	// Calculate total hours
	clockInTime, err := time.Parse("15:04:05", clockInStr.String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid clock in time"})
		return
	}

	clockOutTime, err := time.Parse("15:04:05", now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid clock out time"})
		return
	}

	// Calculate total hours (considering break time - assume 1 hour break if more than 6 hours)
	duration := clockOutTime.Sub(clockInTime)
	totalHours := duration.Hours()

	// Deduct break time if worked more than 6 hours
	if totalHours > 6 {
		totalHours -= 1 // 1 hour break
	}

	// Calculate overtime (over 8 hours)
	overtimeHours := 0.0
	if totalHours > 8 {
		overtimeHours = totalHours - 8
	}

	// Update attendance record
	_, err = database.DB.Exec(`
		UPDATE attendance_logs 
		SET clock_out = ?, total_hours = ?, overtime_hours = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, now, totalHours, overtimeHours, attendanceID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clock out"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Clocked out successfully",
		"time":           now,
		"total_hours":    fmt.Sprintf("%.2f", totalHours),
		"overtime_hours": fmt.Sprintf("%.2f", overtimeHours),
	})
}