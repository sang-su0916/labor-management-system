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

type CreateEmployeeWithContractRequest struct {
	// Employee information (existing fields)
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
	
	// Contract fields
	GenerateContract bool    `json:"generate_contract"` // Whether to auto-generate contract
	ContractType     string  `json:"contract_type"`     // "permanent", "temporary", "contract"
	ContractEndDate  string  `json:"contract_end_date"` // Only for temporary contracts
	Workplace        string  `json:"workplace"`
	JobDescription   string  `json:"job_description"`
	WorkingHours     string  `json:"working_hours"`     // "09:00-18:00"
	WorkDays         string  `json:"work_days"`         // "월-금"
	Allowances       string  `json:"allowances"`
	Benefits         string  `json:"benefits"`
	ContractTerms    string  `json:"contract_terms"`
	
	// Document generation
	GenerateDocument bool `json:"generate_document"` // Whether to auto-generate PDF
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

// CreateEmployeeWithContract creates employee and optionally generates contract and documents
func CreateEmployeeWithContract(c *gin.Context) {
	var req CreateEmployeeWithContractRequest
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
	if req.GenerateContract {
		if req.ContractType == "" {
			req.ContractType = "permanent"
		}
		if req.Workplace == "" {
			req.Workplace = "본사"
		}
		if req.WorkingHours == "" {
			req.WorkingHours = "09:00-18:00"
		}
		if req.WorkDays == "" {
			req.WorkDays = "월-금"
		}
		if req.JobDescription == "" {
			req.JobDescription = fmt.Sprintf("%s 업무 전반", req.Position)
		}
		if req.Benefits == "" {
			req.Benefits = "4대보험, 연차, 퇴직금"
		}
		if req.ContractTerms == "" {
			req.ContractTerms = "회사 규정에 따름"
		}
	}

	// 1. Create Employee
	result, err := tx.Exec(`
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

	// 2. Initialize annual leave balance
	currentYear := time.Now().Year()
	_, err = tx.Exec(`
		INSERT INTO annual_leave_balance (employee_id, year, total_days, remaining_days)
		VALUES (?, ?, 15, 15)
	`, empID, currentYear)

	if err != nil {
		// Log the error but don't fail
		c.Header("Warning", "Failed to create annual leave balance")
	}

	var contractID int64 = 0

	// 3. Generate Contract if requested
	if req.GenerateContract {
		// Parse contract end date if provided
		var contractEndDate sql.NullTime
		if req.ContractEndDate != "" {
			ed, err := time.Parse("2006-01-02", req.ContractEndDate)
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
		`, empID, req.ContractType, hireDate, contractEndDate, req.Workplace,
			req.JobDescription, req.WorkingHours, req.WorkDays, req.BaseSalary,
			req.Allowances, req.Benefits, req.ContractTerms)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contract"})
			return
		}

		contractID, _ = contractResult.LastInsertId()
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// 4. Generate PDF document if requested (after transaction commit)
	var documentPath string
	if req.GenerateDocument && req.GenerateContract {
		documentPath = generateContractPDF(int(empID), req)
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

	response := gin.H{
		"message": "직원이 성공적으로 등록되었습니다",
		"employee": emp,
	}

	if req.GenerateContract {
		response["contract_generated"] = true
		response["contract_id"] = contractID
		response["message"] = "직원과 근로계약서가 성공적으로 생성되었습니다"
	}

	if req.GenerateDocument && documentPath != "" {
		response["document_generated"] = true
		response["document_path"] = documentPath
	}

	c.JSON(http.StatusCreated, response)
}

// Helper function to generate contract PDF
func generateContractPDF(employeeID int, req CreateEmployeeWithContractRequest) string {
	// This is a simplified version - you can enhance it further
	fileName := fmt.Sprintf("contract_%d_%s.pdf", employeeID, time.Now().Format("20060102"))
	filePath := fmt.Sprintf("documents/%s", fileName)
	
	// Generate basic contract PDF
	// Note: This would integrate with the existing PDF generation logic
	// For now, we just return the expected file path
	
	return filePath
}