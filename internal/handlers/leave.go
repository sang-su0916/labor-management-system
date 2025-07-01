package handlers

import (
	"database/sql"
	"labor-management-system/database"
	"labor-management-system/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateLeaveRequestBody struct {
	EmployeeID    int     `json:"employee_id" binding:"required"`
	LeaveType     string  `json:"leave_type" binding:"required"`
	StartDate     string  `json:"start_date" binding:"required"`
	EndDate       string  `json:"end_date" binding:"required"`
	DaysRequested float64 `json:"days_requested" binding:"required"`
	Reason        string  `json:"reason"`
}

type ApproveLeaveRequestBody struct {
	ApprovedBy int `json:"approved_by" binding:"required"`
}

type RejectLeaveRequestBody struct {
	ApprovedBy      int    `json:"approved_by" binding:"required"`
	RejectionReason string `json:"rejection_reason" binding:"required"`
}

func GetLeaveRequests(c *gin.Context) {
	// Get query parameters for filtering
	status := c.Query("status")
	employeeID := c.Query("employee_id")

	query := `
		SELECT l.id, l.employee_id, l.leave_type, l.start_date, l.end_date, 
		       l.days_requested, l.reason, l.status, l.approved_by, l.approved_at, 
		       l.rejection_reason, l.created_at, l.updated_at,
		       e.name as employee_name, e.employee_number,
		       u.username as approved_by_name
		FROM leave_requests l
		JOIN employees e ON l.employee_id = e.id
		LEFT JOIN users u ON l.approved_by = u.id
		WHERE 1=1
	`

	args := []interface{}{}

	if status != "" {
		query += " AND l.status = ?"
		args = append(args, status)
	}

	if employeeID != "" {
		query += " AND l.employee_id = ?"
		args = append(args, employeeID)
	}

	query += " ORDER BY l.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var leaves []map[string]interface{}
	for rows.Next() {
		var leave models.LeaveRequest
		var employeeName, employeeNumber string
		var approvedByName sql.NullString

		err := rows.Scan(
			&leave.ID, &leave.EmployeeID, &leave.LeaveType, &leave.StartDate,
			&leave.EndDate, &leave.DaysRequested, &leave.Reason, &leave.Status,
			&leave.ApprovedBy, &leave.ApprovedAt, &leave.RejectionReason,
			&leave.CreatedAt, &leave.UpdatedAt, &employeeName, &employeeNumber,
			&approvedByName,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan leave request"})
			return
		}

		leaveData := map[string]interface{}{
			"leave":           leave,
			"employee_name":   employeeName,
			"employee_number": employeeNumber,
		}

		if approvedByName.Valid {
			leaveData["approved_by_name"] = approvedByName.String
		}

		leaves = append(leaves, leaveData)
	}

	c.JSON(http.StatusOK, gin.H{"leaves": leaves})
}

func GetLeaveRequest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave request ID"})
		return
	}

	var leave models.LeaveRequest
	var employeeName, employeeNumber string
	var approvedByName sql.NullString

	err = database.DB.QueryRow(`
		SELECT l.id, l.employee_id, l.leave_type, l.start_date, l.end_date, 
		       l.days_requested, l.reason, l.status, l.approved_by, l.approved_at, 
		       l.rejection_reason, l.created_at, l.updated_at,
		       e.name as employee_name, e.employee_number,
		       u.username as approved_by_name
		FROM leave_requests l
		JOIN employees e ON l.employee_id = e.id
		LEFT JOIN users u ON l.approved_by = u.id
		WHERE l.id = ?
	`, id).Scan(
		&leave.ID, &leave.EmployeeID, &leave.LeaveType, &leave.StartDate,
		&leave.EndDate, &leave.DaysRequested, &leave.Reason, &leave.Status,
		&leave.ApprovedBy, &leave.ApprovedAt, &leave.RejectionReason,
		&leave.CreatedAt, &leave.UpdatedAt, &employeeName, &employeeNumber,
		&approvedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Leave request not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	leaveData := map[string]interface{}{
		"leave":           leave,
		"employee_name":   employeeName,
		"employee_number": employeeNumber,
	}

	if approvedByName.Valid {
		leaveData["approved_by_name"] = approvedByName.String
	}

	c.JSON(http.StatusOK, leaveData)
}

func CreateLeaveRequest(c *gin.Context) {
	var req CreateLeaveRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format (YYYY-MM-DD)"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format (YYYY-MM-DD)"})
		return
	}

	// Validate dates
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date cannot be before start date"})
		return
	}

	// Check annual leave balance if leave type is annual
	if req.LeaveType == "annual" {
		var remainingDays float64
		currentYear := time.Now().Year()
		
		err = database.DB.QueryRow(
			"SELECT remaining_days FROM annual_leave_balance WHERE employee_id = ? AND year = ?",
			req.EmployeeID, currentYear,
		).Scan(&remainingDays)

		if err != nil {
			if err == sql.ErrNoRows {
				// Create annual leave balance record if doesn't exist
				_, err = database.DB.Exec(`
					INSERT INTO annual_leave_balance (employee_id, year, total_days, remaining_days)
					VALUES (?, ?, 15, 15)
				`, req.EmployeeID, currentYear)
				
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create annual leave balance"})
					return
				}
				remainingDays = 15
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}
		}

		if req.DaysRequested > remainingDays {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficient annual leave balance",
				"remaining_days": remainingDays,
			})
			return
		}
	}

	// Insert leave request
	result, err := database.DB.Exec(`
		INSERT INTO leave_requests (employee_id, leave_type, start_date, end_date, 
		                           days_requested, reason)
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.EmployeeID, req.LeaveType, startDate, endDate, req.DaysRequested, req.Reason)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create leave request"})
		return
	}

	leaveID, _ := result.LastInsertId()

	// Retrieve created leave request with employee info
	var leave models.LeaveRequest
	var employeeName, employeeNumber string

	err = database.DB.QueryRow(`
		SELECT l.id, l.employee_id, l.leave_type, l.start_date, l.end_date, 
		       l.days_requested, l.reason, l.status, l.approved_by, l.approved_at, 
		       l.rejection_reason, l.created_at, l.updated_at,
		       e.name as employee_name, e.employee_number
		FROM leave_requests l
		JOIN employees e ON l.employee_id = e.id
		WHERE l.id = ?
	`, leaveID).Scan(
		&leave.ID, &leave.EmployeeID, &leave.LeaveType, &leave.StartDate,
		&leave.EndDate, &leave.DaysRequested, &leave.Reason, &leave.Status,
		&leave.ApprovedBy, &leave.ApprovedAt, &leave.RejectionReason,
		&leave.CreatedAt, &leave.UpdatedAt, &employeeName, &employeeNumber,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created leave request"})
		return
	}

	leaveData := map[string]interface{}{
		"leave":           leave,
		"employee_name":   employeeName,
		"employee_number": employeeNumber,
	}

	c.JSON(http.StatusCreated, leaveData)
}

func ApproveLeaveRequest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave request ID"})
		return
	}

	var req ApproveLeaveRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get leave request details
	var leave models.LeaveRequest
	err = database.DB.QueryRow(`
		SELECT id, employee_id, leave_type, start_date, end_date, days_requested, status
		FROM leave_requests WHERE id = ?
	`, id).Scan(
		&leave.ID, &leave.EmployeeID, &leave.LeaveType, &leave.StartDate,
		&leave.EndDate, &leave.DaysRequested, &leave.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Leave request not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if leave.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Leave request is not pending"})
		return
	}

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Update leave request status
	_, err = tx.Exec(`
		UPDATE leave_requests 
		SET status = 'approved', approved_by = ?, approved_at = CURRENT_TIMESTAMP, 
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, req.ApprovedBy, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve leave request"})
		return
	}

	// Update annual leave balance if it's annual leave
	if leave.LeaveType == "annual" {
		currentYear := time.Now().Year()
		_, err = tx.Exec(`
			UPDATE annual_leave_balance 
			SET used_days = used_days + ?, remaining_days = remaining_days - ?, 
			    updated_at = CURRENT_TIMESTAMP
			WHERE employee_id = ? AND year = ?
		`, leave.DaysRequested, leave.DaysRequested, leave.EmployeeID, currentYear)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update annual leave balance"})
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Leave request approved successfully"})
}

func RejectLeaveRequest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave request ID"})
		return
	}

	var req RejectLeaveRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if leave request exists and is pending
	var status string
	err = database.DB.QueryRow("SELECT status FROM leave_requests WHERE id = ?", id).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Leave request not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Leave request is not pending"})
		return
	}

	// Update leave request status
	_, err = database.DB.Exec(`
		UPDATE leave_requests 
		SET status = 'rejected', approved_by = ?, approved_at = CURRENT_TIMESTAMP, 
		    rejection_reason = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, req.ApprovedBy, req.RejectionReason, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject leave request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Leave request rejected successfully"})
}