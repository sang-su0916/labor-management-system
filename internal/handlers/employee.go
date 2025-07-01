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

type CreateEmployeeRequest struct {
	EmployeeNumber string    `json:"employee_number" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	NameEn         string    `json:"name_en"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	Address        string    `json:"address"`
	BirthDate      string    `json:"birth_date"`
	HireDate       string    `json:"hire_date" binding:"required"`
	Department     string    `json:"department"`
	Position       string    `json:"position"`
	EmploymentType string    `json:"employment_type"`
	SalaryType     string    `json:"salary_type"`
	BaseSalary     float64   `json:"base_salary"`
}

func GetEmployees(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, user_id, employee_number, name, name_en, phone, email, address, 
		       birth_date, hire_date, department, position, employment_type, status, 
		       salary_type, base_salary, created_at, updated_at 
		FROM employees 
		WHERE status != 'terminated'
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		err := rows.Scan(
			&emp.ID, &emp.UserID, &emp.EmployeeNumber, &emp.Name, &emp.NameEn,
			&emp.Phone, &emp.Email, &emp.Address, &emp.BirthDate, &emp.HireDate,
			&emp.Department, &emp.Position, &emp.EmploymentType, &emp.Status,
			&emp.SalaryType, &emp.BaseSalary, &emp.CreatedAt, &emp.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan employee"})
			return
		}
		employees = append(employees, emp)
	}

	c.JSON(http.StatusOK, gin.H{"employees": employees})
}

func GetEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var emp models.Employee
	err = database.DB.QueryRow(`
		SELECT id, user_id, employee_number, name, name_en, phone, email, address, 
		       birth_date, hire_date, department, position, employment_type, status, 
		       salary_type, base_salary, created_at, updated_at 
		FROM employees WHERE id = ?
	`, id).Scan(
		&emp.ID, &emp.UserID, &emp.EmployeeNumber, &emp.Name, &emp.NameEn,
		&emp.Phone, &emp.Email, &emp.Address, &emp.BirthDate, &emp.HireDate,
		&emp.Department, &emp.Position, &emp.EmploymentType, &emp.Status,
		&emp.SalaryType, &emp.BaseSalary, &emp.CreatedAt, &emp.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"employee": emp})
}

func CreateEmployee(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse hire date
	hireDate, err := time.Parse("2006-01-02", req.HireDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hire date format (YYYY-MM-DD)"})
		return
	}

	// Parse birth date if provided
	var birthDate sql.NullTime
	if req.BirthDate != "" {
		bd, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birth date format (YYYY-MM-DD)"})
			return
		}
		birthDate = sql.NullTime{Time: bd, Valid: true}
	}

	// Set defaults
	if req.EmploymentType == "" {
		req.EmploymentType = "regular"
	}
	if req.SalaryType == "" {
		req.SalaryType = "monthly"
	}

	// Insert employee
	result, err := database.DB.Exec(`
		INSERT INTO employees (employee_number, name, name_en, phone, email, address, 
		                      birth_date, hire_date, department, position, employment_type, 
		                      salary_type, base_salary)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.EmployeeNumber, req.Name, req.NameEn, req.Phone, req.Email, req.Address,
		birthDate, hireDate, req.Department, req.Position, req.EmploymentType,
		req.SalaryType, req.BaseSalary)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Employee number already exists"})
		return
	}

	empID, _ := result.LastInsertId()

	// Initialize annual leave balance for the employee
	currentYear := time.Now().Year()
	_, err = database.DB.Exec(`
		INSERT INTO annual_leave_balance (employee_id, year, total_days, remaining_days)
		VALUES (?, ?, 15, 15)
	`, empID, currentYear)

	if err != nil {
		// Log the error but don't fail the employee creation
		// log.Printf("Failed to create annual leave balance: %v", err)
	}

	// Retrieve created employee
	var emp models.Employee
	err = database.DB.QueryRow(`
		SELECT id, user_id, employee_number, name, name_en, phone, email, address, 
		       birth_date, hire_date, department, position, employment_type, status, 
		       salary_type, base_salary, created_at, updated_at 
		FROM employees WHERE id = ?
	`, empID).Scan(
		&emp.ID, &emp.UserID, &emp.EmployeeNumber, &emp.Name, &emp.NameEn,
		&emp.Phone, &emp.Email, &emp.Address, &emp.BirthDate, &emp.HireDate,
		&emp.Department, &emp.Position, &emp.EmploymentType, &emp.Status,
		&emp.SalaryType, &emp.BaseSalary, &emp.CreatedAt, &emp.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created employee"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"employee": emp})
}

func UpdateEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse hire date
	hireDate, err := time.Parse("2006-01-02", req.HireDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hire date format (YYYY-MM-DD)"})
		return
	}

	// Parse birth date if provided
	var birthDate sql.NullTime
	if req.BirthDate != "" {
		bd, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birth date format (YYYY-MM-DD)"})
			return
		}
		birthDate = sql.NullTime{Time: bd, Valid: true}
	}

	// Update employee
	_, err = database.DB.Exec(`
		UPDATE employees SET name = ?, name_en = ?, phone = ?, email = ?, address = ?, 
		                    birth_date = ?, hire_date = ?, department = ?, position = ?, 
		                    employment_type = ?, salary_type = ?, base_salary = ?, 
		                    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, req.Name, req.NameEn, req.Phone, req.Email, req.Address, birthDate, hireDate,
		req.Department, req.Position, req.EmploymentType, req.SalaryType, req.BaseSalary, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}

	// Retrieve updated employee
	var emp models.Employee
	err = database.DB.QueryRow(`
		SELECT id, user_id, employee_number, name, name_en, phone, email, address, 
		       birth_date, hire_date, department, position, employment_type, status, 
		       salary_type, base_salary, created_at, updated_at 
		FROM employees WHERE id = ?
	`, id).Scan(
		&emp.ID, &emp.UserID, &emp.EmployeeNumber, &emp.Name, &emp.NameEn,
		&emp.Phone, &emp.Email, &emp.Address, &emp.BirthDate, &emp.HireDate,
		&emp.Department, &emp.Position, &emp.EmploymentType, &emp.Status,
		&emp.SalaryType, &emp.BaseSalary, &emp.CreatedAt, &emp.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"employee": emp})
}

func DeleteEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	// Soft delete by updating status
	_, err = database.DB.Exec(`
		UPDATE employees SET status = 'terminated', updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}