package handlers

import (
	"database/sql"
	"fmt"
	"labor-management-system/database"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

type DocumentTemplate struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Variables string `json:"variables"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GeneratedDocument struct {
	ID           int    `json:"id"`
	EmployeeID   *int   `json:"employee_id"`
	TemplateID   int    `json:"template_id"`
	DocumentType string `json:"document_type"`
	FilePath     string `json:"file_path"`
	GeneratedBy  int    `json:"generated_by"`
	GeneratedAt  string `json:"generated_at"`
}

func GetDocumentTemplates(c *gin.Context) {
	db := database.GetDB()
	
	rows, err := db.Query(`
		SELECT id, name, type, content, COALESCE(variables, ''), is_active, 
		       created_at, updated_at
		FROM document_templates 
		WHERE is_active = 1
		ORDER BY type, name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch templates"})
		return
	}
	defer rows.Close()

	var templates []DocumentTemplate
	for rows.Next() {
		var template DocumentTemplate
		err := rows.Scan(&template.ID, &template.Name, &template.Type, 
			&template.Content, &template.Variables, &template.IsActive,
			&template.CreatedAt, &template.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan template"})
			return
		}
		templates = append(templates, template)
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

func GenerateDocument(c *gin.Context) {
	docType := c.Param("type")
	employeeIDStr := c.Query("employee_id")
	
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var employeeID *int
	if employeeIDStr != "" {
		if id, err := strconv.Atoi(employeeIDStr); err == nil {
			employeeID = &id
		}
	}

	switch docType {
	case "payslip":
		generatePayslip(c, employeeID, userID.(int))
	case "employment_certificate":
		generateEmploymentCertificate(c, employeeID, userID.(int))
	case "contract":
		generateContract(c, employeeID, userID.(int))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported document type"})
	}
}

func generatePayslip(c *gin.Context, employeeID *int, userID int) {
	if employeeID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employee ID is required"})
		return
	}

	db := database.GetDB()
	
	// Get employee info
	var employee struct {
		Name       string
		Department string
		Position   string
	}
	
	err := db.QueryRow(`
		SELECT name, department, position 
		FROM employees 
		WHERE id = ?
	`, *employeeID).Scan(&employee.Name, &employee.Department, &employee.Position)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// Get latest payroll record
	var payroll struct {
		PayPeriod       string
		BaseSalary      float64
		Allowances      float64
		Bonus           float64
		GrossPay        float64
		IncomeTax       float64
		LocalTax        float64
		NationalPension float64
		HealthInsurance float64
		EmploymentIns   float64
		LongTermCare    float64
		OtherDeductions float64
		TotalDeductions float64
		NetPay          float64
	}
	
	err = db.QueryRow(`
		SELECT pay_period, base_salary, allowances, bonus, gross_pay,
		       income_tax, local_tax, national_pension, health_insurance,
		       employment_insurance, long_term_care, other_deductions,
		       total_deductions, net_pay
		FROM payroll_records 
		WHERE employee_id = ?
		ORDER BY created_at DESC LIMIT 1
	`, *employeeID).Scan(&payroll.PayPeriod, &payroll.BaseSalary, &payroll.Allowances,
		&payroll.Bonus, &payroll.GrossPay, &payroll.IncomeTax, &payroll.LocalTax,
		&payroll.NationalPension, &payroll.HealthInsurance, &payroll.EmploymentIns,
		&payroll.LongTermCare, &payroll.OtherDeductions, &payroll.TotalDeductions,
		&payroll.NetPay)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No payroll record found"})
		return
	}

	// Generate PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	
	// Title
	pdf.Cell(40, 10, "급여명세서")
	pdf.Ln(20)
	
	// Employee info
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("성명: %s", employee.Name))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("부서: %s", employee.Department))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("직급: %s", employee.Position))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("급여기간: %s", payroll.PayPeriod))
	pdf.Ln(15)
	
	// Payroll details
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "지급내역")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(40, 8, fmt.Sprintf("기본급: %,.0f원", payroll.BaseSalary))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("수당: %,.0f원", payroll.Allowances))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("상여금: %,.0f원", payroll.Bonus))
	pdf.Ln(6)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(40, 8, fmt.Sprintf("총 지급액: %,.0f원", payroll.GrossPay))
	pdf.Ln(15)
	
	// Deductions
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "공제내역")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(40, 8, fmt.Sprintf("소득세: %,.0f원", payroll.IncomeTax))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("지방소득세: %,.0f원", payroll.LocalTax))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("국민연금: %,.0f원", payroll.NationalPension))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("건강보험: %,.0f원", payroll.HealthInsurance))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("고용보험: %,.0f원", payroll.EmploymentIns))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("장기요양보험: %,.0f원", payroll.LongTermCare))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("기타공제: %,.0f원", payroll.OtherDeductions))
	pdf.Ln(6)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(40, 8, fmt.Sprintf("총 공제액: %,.0f원", payroll.TotalDeductions))
	pdf.Ln(15)
	
	// Net pay
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, fmt.Sprintf("실지급액: %,.0f원", payroll.NetPay))
	
	// Save PDF
	fileName := fmt.Sprintf("payslip_%d_%s.pdf", *employeeID, time.Now().Format("20060102"))
	filePath := filepath.Join("documents", fileName)
	
	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}
	
	// Record generation
	_, err = db.Exec(`
		INSERT INTO generated_documents (employee_id, template_id, document_type, file_path, generated_by)
		VALUES (?, 1, 'payslip', ?, ?)
	`, *employeeID, filePath, userID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record document generation"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Payslip generated successfully",
		"file_path": filePath,
	})
}

func generateEmploymentCertificate(c *gin.Context, employeeID *int, userID int) {
	if employeeID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employee ID is required"})
		return
	}

	db := database.GetDB()
	
	var employee struct {
		Name       string
		Department string
		Position   string
		HireDate   string
	}
	
	err := db.QueryRow(`
		SELECT name, department, position, hire_date
		FROM employees 
		WHERE id = ?
	`, *employeeID).Scan(&employee.Name, &employee.Department, &employee.Position, &employee.HireDate)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// Generate PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 18)
	
	// Title
	pdf.Cell(40, 20, "재직증명서")
	pdf.Ln(30)
	
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("성    명: %s", employee.Name))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("부    서: %s", employee.Department))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("직    급: %s", employee.Position))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("입사일자: %s", employee.HireDate))
	pdf.Ln(20)
	
	pdf.Cell(40, 10, "위 사람은 현재 본 회사에 재직 중임을 증명합니다.")
	pdf.Ln(30)
	
	pdf.Cell(40, 10, fmt.Sprintf("발급일자: %s", time.Now().Format("2006년 01월 02일")))
	pdf.Ln(20)
	
	pdf.Cell(40, 10, "발급기관: (주)노무관리시스템")
	
	fileName := fmt.Sprintf("certificate_%d_%s.pdf", *employeeID, time.Now().Format("20060102"))
	filePath := filepath.Join("documents", fileName)
	
	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}
	
	// Record generation
	_, err = db.Exec(`
		INSERT INTO generated_documents (employee_id, template_id, document_type, file_path, generated_by)
		VALUES (?, 2, 'employment_certificate', ?, ?)
	`, *employeeID, filePath, userID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record document generation"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Employment certificate generated successfully",
		"file_path": filePath,
	})
}

func generateContract(c *gin.Context, employeeID *int, userID int) {
	if employeeID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employee ID is required"})
		return
	}

	db := database.GetDB()
	
	var contract struct {
		EmployeeName string
		Position     string
		Department   string
		StartDate    string
		EndDate      sql.NullString
		Salary       float64
		WorkingHours int
		WorkDays     string
	}
	
	err := db.QueryRow(`
		SELECT e.name, ec.position, e.department, ec.start_date, ec.end_date,
		       ec.salary, ec.working_hours, ec.work_days
		FROM employment_contracts ec
		JOIN employees e ON ec.employee_id = e.id
		WHERE ec.employee_id = ? AND ec.status = 'active'
		ORDER BY ec.created_at DESC LIMIT 1
	`, *employeeID).Scan(&contract.EmployeeName, &contract.Position, &contract.Department,
		&contract.StartDate, &contract.EndDate, &contract.Salary, &contract.WorkingHours,
		&contract.WorkDays)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active contract found"})
		return
	}

	// Generate PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	
	pdf.Cell(40, 20, "근로계약서")
	pdf.Ln(30)
	
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("성    명: %s", contract.EmployeeName))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("부    서: %s", contract.Department))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("직    급: %s", contract.Position))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("계약기간: %s", contract.StartDate))
	if contract.EndDate.Valid {
		pdf.Cell(40, 10, fmt.Sprintf(" ~ %s", contract.EndDate.String))
	} else {
		pdf.Cell(40, 10, " ~ 무기한")
	}
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("급    여: %,.0f원", contract.Salary))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("근무시간: %d시간", contract.WorkingHours))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("근무요일: %s", contract.WorkDays))
	
	fileName := fmt.Sprintf("contract_%d_%s.pdf", *employeeID, time.Now().Format("20060102"))
	filePath := filepath.Join("documents", fileName)
	
	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}
	
	// Record generation
	_, err = db.Exec(`
		INSERT INTO generated_documents (employee_id, template_id, document_type, file_path, generated_by)
		VALUES (?, 3, 'contract', ?, ?)
	`, *employeeID, filePath, userID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record document generation"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Contract generated successfully",
		"file_path": filePath,
	})
}

func GetEmployeeDocuments(c *gin.Context) {
	employeeIDStr := c.Param("id")
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	db := database.GetDB()
	
	rows, err := db.Query(`
		SELECT gd.id, gd.document_type, gd.file_path, gd.generated_at,
		       u.username as generated_by_username
		FROM generated_documents gd
		JOIN users u ON gd.generated_by = u.id
		WHERE gd.employee_id = ?
		ORDER BY gd.generated_at DESC
	`, employeeID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}
	defer rows.Close()

	var documents []map[string]interface{}
	for rows.Next() {
		var doc map[string]interface{} = make(map[string]interface{})
		var id int
		var docType, filePath, generatedAt, generatedBy string
		
		err := rows.Scan(&id, &docType, &filePath, &generatedAt, &generatedBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan document"})
			return
		}
		
		doc["id"] = id
		doc["document_type"] = docType
		doc["file_path"] = filePath
		doc["generated_at"] = generatedAt
		doc["generated_by"] = generatedBy
		
		documents = append(documents, doc)
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}