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

type CreateContractRequest struct {
	EmployeeID     int     `json:"employee_id" binding:"required"`
	ContractType   string  `json:"contract_type" binding:"required"`
	StartDate      string  `json:"start_date" binding:"required"`
	EndDate        string  `json:"end_date"`
	Workplace      string  `json:"workplace" binding:"required"`
	JobDescription string  `json:"job_description"`
	WorkingHours   string  `json:"working_hours" binding:"required"`
	WorkDays       string  `json:"work_days" binding:"required"`
	BaseSalary     float64 `json:"base_salary" binding:"required"`
	Allowances     string  `json:"allowances"`
	Benefits       string  `json:"benefits"`
	ContractTerms  string  `json:"contract_terms"`
}

// Request struct for creating contract with new employee
type CreateContractWithEmployeeRequest struct {
	// Employee information
	EmployeeName     string `json:"employee_name" binding:"required"`
	EmployeeNumber   string `json:"employee_number" binding:"required"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	Address          string `json:"address"`
	BirthDate        string `json:"birth_date"`
	Department       string `json:"department" binding:"required"`
	Position         string `json:"position" binding:"required"`
	EmploymentType   string `json:"employment_type"`
	SalaryType       string `json:"salary_type"`
	
	// Contract information
	ContractType   string  `json:"contract_type" binding:"required"`
	StartDate      string  `json:"start_date" binding:"required"`
	EndDate        string  `json:"end_date"`
	Workplace      string  `json:"workplace" binding:"required"`
	JobDescription string  `json:"job_description"`
	WorkingHours   string  `json:"working_hours" binding:"required"`
	WorkDays       string  `json:"work_days" binding:"required"`
	BaseSalary     float64 `json:"base_salary" binding:"required"`
	Allowances     string  `json:"allowances"`
	Benefits       string  `json:"benefits"`
	ContractTerms  string  `json:"contract_terms"`
}

func GetContracts(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT c.id, c.employee_id, c.contract_type, c.start_date, c.end_date, 
		       c.workplace, c.job_description, c.working_hours, c.work_days, 
		       c.base_salary, c.allowances, c.benefits, c.contract_terms, 
		       c.signed_date, c.is_active, c.created_at, c.updated_at,
		       e.name as employee_name, e.employee_number
		FROM employment_contracts c
		JOIN employees e ON c.employee_id = e.id
		WHERE c.is_active = 1
		ORDER BY c.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var contracts []map[string]interface{}
	for rows.Next() {
		var contract models.EmploymentContract
		var employeeName, employeeNumber string
		
		err := rows.Scan(
			&contract.ID, &contract.EmployeeID, &contract.ContractType, 
			&contract.StartDate, &contract.EndDate, &contract.Workplace,
			&contract.JobDescription, &contract.WorkingHours, &contract.WorkDays,
			&contract.BaseSalary, &contract.Allowances, &contract.Benefits,
			&contract.ContractTerms, &contract.SignedDate, &contract.IsActive,
			&contract.CreatedAt, &contract.UpdatedAt, &employeeName, &employeeNumber,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan contract"})
			return
		}

		contractData := map[string]interface{}{
			"contract":        contract,
			"employee_name":   employeeName,
			"employee_number": employeeNumber,
		}
		contracts = append(contracts, contractData)
	}

	c.JSON(http.StatusOK, gin.H{"contracts": contracts})
}

func GetContract(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var contract models.EmploymentContract
	var employeeName, employeeNumber string

	err = database.DB.QueryRow(`
		SELECT c.id, c.employee_id, c.contract_type, c.start_date, c.end_date, 
		       c.workplace, c.job_description, c.working_hours, c.work_days, 
		       c.base_salary, c.allowances, c.benefits, c.contract_terms, 
		       c.signed_date, c.is_active, c.created_at, c.updated_at,
		       e.name as employee_name, e.employee_number
		FROM employment_contracts c
		JOIN employees e ON c.employee_id = e.id
		WHERE c.id = ?
	`, id).Scan(
		&contract.ID, &contract.EmployeeID, &contract.ContractType,
		&contract.StartDate, &contract.EndDate, &contract.Workplace,
		&contract.JobDescription, &contract.WorkingHours, &contract.WorkDays,
		&contract.BaseSalary, &contract.Allowances, &contract.Benefits,
		&contract.ContractTerms, &contract.SignedDate, &contract.IsActive,
		&contract.CreatedAt, &contract.UpdatedAt, &employeeName, &employeeNumber,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	contractData := map[string]interface{}{
		"contract":        contract,
		"employee_name":   employeeName,
		"employee_number": employeeNumber,
	}

	c.JSON(http.StatusOK, contractData)
}

func CreateContract(c *gin.Context) {
	var req CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format (YYYY-MM-DD)"})
		return
	}

	// Parse end date if provided
	var endDate sql.NullTime
	if req.EndDate != "" {
		ed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format (YYYY-MM-DD)"})
			return
		}
		endDate = sql.NullTime{Time: ed, Valid: true}
	}

	// Deactivate existing contracts for the employee
	_, err = database.DB.Exec(`
		UPDATE employment_contracts SET is_active = 0, updated_at = CURRENT_TIMESTAMP 
		WHERE employee_id = ? AND is_active = 1
	`, req.EmployeeID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate existing contracts"})
		return
	}

	// Insert new contract
	result, err := database.DB.Exec(`
		INSERT INTO employment_contracts (employee_id, contract_type, start_date, end_date, 
		                                 workplace, job_description, working_hours, work_days, 
		                                 base_salary, allowances, benefits, contract_terms)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.EmployeeID, req.ContractType, startDate, endDate, req.Workplace,
		req.JobDescription, req.WorkingHours, req.WorkDays, req.BaseSalary,
		req.Allowances, req.Benefits, req.ContractTerms)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contract"})
		return
	}

	contractID, _ := result.LastInsertId()

	// Retrieve created contract with employee info
	var contract models.EmploymentContract
	var employeeName, employeeNumber string

	err = database.DB.QueryRow(`
		SELECT c.id, c.employee_id, c.contract_type, c.start_date, c.end_date, 
		       c.workplace, c.job_description, c.working_hours, c.work_days, 
		       c.base_salary, c.allowances, c.benefits, c.contract_terms, 
		       c.signed_date, c.is_active, c.created_at, c.updated_at,
		       e.name as employee_name, e.employee_number
		FROM employment_contracts c
		JOIN employees e ON c.employee_id = e.id
		WHERE c.id = ?
	`, contractID).Scan(
		&contract.ID, &contract.EmployeeID, &contract.ContractType,
		&contract.StartDate, &contract.EndDate, &contract.Workplace,
		&contract.JobDescription, &contract.WorkingHours, &contract.WorkDays,
		&contract.BaseSalary, &contract.Allowances, &contract.Benefits,
		&contract.ContractTerms, &contract.SignedDate, &contract.IsActive,
		&contract.CreatedAt, &contract.UpdatedAt, &employeeName, &employeeNumber,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created contract"})
		return
	}

	contractData := map[string]interface{}{
		"contract":        contract,
		"employee_name":   employeeName,
		"employee_number": employeeNumber,
	}

	c.JSON(http.StatusCreated, contractData)
}

func UpdateContract(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var req CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format (YYYY-MM-DD)"})
		return
	}

	// Parse end date if provided
	var endDate sql.NullTime
	if req.EndDate != "" {
		ed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format (YYYY-MM-DD)"})
			return
		}
		endDate = sql.NullTime{Time: ed, Valid: true}
	}

	// Update contract
	_, err = database.DB.Exec(`
		UPDATE employment_contracts SET contract_type = ?, start_date = ?, end_date = ?, 
		                               workplace = ?, job_description = ?, working_hours = ?, 
		                               work_days = ?, base_salary = ?, allowances = ?, 
		                               benefits = ?, contract_terms = ?, 
		                               updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, req.ContractType, startDate, endDate, req.Workplace, req.JobDescription,
		req.WorkingHours, req.WorkDays, req.BaseSalary, req.Allowances,
		req.Benefits, req.ContractTerms, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contract"})
		return
	}

	// Retrieve updated contract with employee info
	var contract models.EmploymentContract
	var employeeName, employeeNumber string

	err = database.DB.QueryRow(`
		SELECT c.id, c.employee_id, c.contract_type, c.start_date, c.end_date, 
		       c.workplace, c.job_description, c.working_hours, c.work_days, 
		       c.base_salary, c.allowances, c.benefits, c.contract_terms, 
		       c.signed_date, c.is_active, c.created_at, c.updated_at,
		       e.name as employee_name, e.employee_number
		FROM employment_contracts c
		JOIN employees e ON c.employee_id = e.id
		WHERE c.id = ?
	`, id).Scan(
		&contract.ID, &contract.EmployeeID, &contract.ContractType,
		&contract.StartDate, &contract.EndDate, &contract.Workplace,
		&contract.JobDescription, &contract.WorkingHours, &contract.WorkDays,
		&contract.BaseSalary, &contract.Allowances, &contract.Benefits,
		&contract.ContractTerms, &contract.SignedDate, &contract.IsActive,
		&contract.CreatedAt, &contract.UpdatedAt, &employeeName, &employeeNumber,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated contract"})
		return
	}

	contractData := map[string]interface{}{
		"contract":        contract,
		"employee_name":   employeeName,
		"employee_number": employeeNumber,
	}

	c.JSON(http.StatusOK, contractData)
}

func DeleteContract(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	// Soft delete by deactivating
	_, err = database.DB.Exec(`
		UPDATE employment_contracts SET is_active = 0, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contract"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contract deleted successfully"})
}

// CreateContractWithEmployee creates both employee and contract in a single transaction
func CreateContractWithEmployee(c *gin.Context) {
	var req CreateContractWithEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Parse dates
	hireDate, err := time.Parse("2006-01-02", req.StartDate) // Use contract start date as hire date
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format (YYYY-MM-DD)"})
		return
	}

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

	// 1. Create Employee
	employeeResult, err := tx.Exec(`
		INSERT INTO employees (employee_number, name, phone, email, address, 
		                      birth_date, hire_date, department, position, employment_type, 
		                      salary_type, base_salary)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.EmployeeNumber, req.EmployeeName, req.Phone, req.Email, req.Address,
		birthDate, hireDate, req.Department, req.Position, req.EmploymentType,
		req.SalaryType, req.BaseSalary)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Employee number already exists or database error"})
		return
	}

	employeeID, err := employeeResult.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get employee ID"})
		return
	}

	// 2. Initialize annual leave balance
	currentYear := time.Now().Year()
	_, err = tx.Exec(`
		INSERT INTO annual_leave_balance (employee_id, year, total_days, remaining_days)
		VALUES (?, ?, 15, 15)
	`, employeeID, currentYear)

	if err != nil {
		// Log but don't fail the transaction
		c.Header("Warning", "Failed to create annual leave balance")
	}

	// 3. Create Contract
	contractStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract start date format"})
		return
	}

	var contractEndDate sql.NullTime
	if req.EndDate != "" {
		ed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract end date format"})
			return
		}
		contractEndDate = sql.NullTime{Time: ed, Valid: true}
	}

	contractResult, err := tx.Exec(`
		INSERT INTO employment_contracts (employee_id, contract_type, start_date, end_date, 
		                                 workplace, job_description, working_hours, work_days, 
		                                 base_salary, allowances, benefits, contract_terms)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, employeeID, req.ContractType, contractStartDate, contractEndDate, req.Workplace,
		req.JobDescription, req.WorkingHours, req.WorkDays, req.BaseSalary,
		req.Allowances, req.Benefits, req.ContractTerms)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contract"})
		return
	}

	contractID, err := contractResult.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get contract ID"})
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Return created records
	c.JSON(http.StatusCreated, gin.H{
		"message": "직원과 근로계약서가 성공적으로 생성되었습니다",
		"employee_id": employeeID,
		"contract_id": contractID,
		"employee": gin.H{
			"id": employeeID,
			"name": req.EmployeeName,
			"employee_number": req.EmployeeNumber,
			"department": req.Department,
			"position": req.Position,
		},
		"contract": gin.H{
			"id": contractID,
			"contract_type": req.ContractType,
			"start_date": req.StartDate,
			"end_date": req.EndDate,
			"base_salary": req.BaseSalary,
		},
	})
}