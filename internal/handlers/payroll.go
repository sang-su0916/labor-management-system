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

type CreatePayrollRequest struct {
	EmployeeID      int     `json:"employee_id" binding:"required"`
	PayPeriodStart  string  `json:"pay_period_start" binding:"required"`
	PayPeriodEnd    string  `json:"pay_period_end" binding:"required"`
	BaseSalary      float64 `json:"base_salary" binding:"required"`
	OvertimeHours   float64 `json:"overtime_hours"`
	HolidayHours    float64 `json:"holiday_hours"`
	Allowances      float64 `json:"allowances"`
	Bonus           float64 `json:"bonus"`
	OtherDeductions float64 `json:"other_deductions"`
}

// PayrollCalculator handles payroll calculations
type PayrollCalculator struct {
	BaseSalary      float64
	OvertimeHours   float64
	HolidayHours    float64
	Allowances      float64
	Bonus           float64
	OtherDeductions float64
}

func (pc *PayrollCalculator) Calculate() map[string]float64 {
	// Korean tax and insurance rates (2024 기준 - 실제 환경에서는 설정에서 관리)
	const (
		overtimeRate         = 1.5  // 연장근로 가산율
		holidayRate          = 2.0  // 휴일근로 가산율
		incomeTaxRate        = 0.05 // 소득세율 (간소화)
		localTaxRate         = 0.1  // 지방소득세율 (소득세의 10%)
		nationalPensionRate  = 0.045 // 국민연금 4.5%
		healthInsuranceRate  = 0.0354 // 건강보험 3.54%
		longTermCareRate     = 0.004564 // 장기요양보험 0.4564%
		employmentInsuranceRate = 0.009 // 고용보험 0.9%
		minWage              = 9860 // 2024년 최저임금 (시급)
	)

	// 기본급 기준 시급 계산 (월급 -> 시급)
	hourlyWage := pc.BaseSalary / (22 * 8) // 월 22일, 일 8시간 기준
	if hourlyWage < minWage {
		hourlyWage = minWage
	}

	// 연장근로수당 계산
	overtimePay := pc.OvertimeHours * hourlyWage * overtimeRate

	// 휴일근로수당 계산
	holidayPay := pc.HolidayHours * hourlyWage * holidayRate

	// 총 지급액 계산
	grossPay := pc.BaseSalary + overtimePay + holidayPay + pc.Allowances + pc.Bonus

	// 4대보험 및 세금 계산
	nationalPension := grossPay * nationalPensionRate
	healthInsurance := grossPay * healthInsuranceRate
	longTermCare := healthInsurance * longTermCareRate
	employmentInsurance := grossPay * employmentInsuranceRate

	// 소득세 계산 (간이세액표 적용 필요하지만 여기서는 단순화)
	incomeTax := grossPay * incomeTaxRate
	localTax := incomeTax * localTaxRate

	// 총 공제액
	totalDeductions := nationalPension + healthInsurance + longTermCare + 
					  employmentInsurance + incomeTax + localTax + pc.OtherDeductions

	// 실지급액
	netPay := grossPay - totalDeductions

	return map[string]float64{
		"overtime_pay":         overtimePay,
		"holiday_pay":          holidayPay,
		"gross_pay":            grossPay,
		"income_tax":           incomeTax,
		"local_tax":            localTax,
		"national_pension":     nationalPension,
		"health_insurance":     healthInsurance,
		"long_term_care":       longTermCare,
		"employment_insurance": employmentInsurance,
		"total_deductions":     totalDeductions,
		"net_pay":              netPay,
	}
}

func GetPayrollRecords(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT p.id, p.employee_id, p.pay_period_start, p.pay_period_end, 
		       p.base_salary, p.overtime_hours, p.overtime_pay, p.holiday_hours, 
		       p.holiday_pay, p.allowances, p.bonus, p.gross_pay, p.income_tax, 
		       p.local_tax, p.national_pension, p.health_insurance, p.employment_insurance, 
		       p.long_term_care, p.other_deductions, p.total_deductions, p.net_pay, 
		       p.pay_date, p.is_paid, p.created_at, p.updated_at,
		       e.name as employee_name, e.employee_number
		FROM payroll_records p
		JOIN employees e ON p.employee_id = e.id
		ORDER BY p.pay_period_start DESC, e.name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var payrolls []map[string]interface{}
	for rows.Next() {
		var payroll models.PayrollRecord
		var employeeName, employeeNumber string
		
		err := rows.Scan(
			&payroll.ID, &payroll.EmployeeID, &payroll.PayPeriodStart, &payroll.PayPeriodEnd,
			&payroll.BaseSalary, &payroll.OvertimeHours, &payroll.OvertimePay, &payroll.HolidayHours,
			&payroll.HolidayPay, &payroll.Allowances, &payroll.Bonus, &payroll.GrossPay,
			&payroll.IncomeTax, &payroll.LocalTax, &payroll.NationalPension, &payroll.HealthInsurance,
			&payroll.EmploymentInsurance, &payroll.LongTermCare, &payroll.OtherDeductions,
			&payroll.TotalDeductions, &payroll.NetPay, &payroll.PayDate, &payroll.IsPaid,
			&payroll.CreatedAt, &payroll.UpdatedAt, &employeeName, &employeeNumber,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan payroll"})
			return
		}

		payrollData := map[string]interface{}{
			"payroll":         payroll,
			"employee_name":   employeeName,
			"employee_number": employeeNumber,
		}
		payrolls = append(payrolls, payrollData)
	}

	c.JSON(http.StatusOK, gin.H{"payrolls": payrolls})
}

func GetPayrollRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return
	}

	var payroll models.PayrollRecord
	var employeeName, employeeNumber string

	err = database.DB.QueryRow(`
		SELECT p.id, p.employee_id, p.pay_period_start, p.pay_period_end, 
		       p.base_salary, p.overtime_hours, p.overtime_pay, p.holiday_hours, 
		       p.holiday_pay, p.allowances, p.bonus, p.gross_pay, p.income_tax, 
		       p.local_tax, p.national_pension, p.health_insurance, p.employment_insurance, 
		       p.long_term_care, p.other_deductions, p.total_deductions, p.net_pay, 
		       p.pay_date, p.is_paid, p.created_at, p.updated_at,
		       e.name as employee_name, e.employee_number
		FROM payroll_records p
		JOIN employees e ON p.employee_id = e.id
		WHERE p.id = ?
	`, id).Scan(
		&payroll.ID, &payroll.EmployeeID, &payroll.PayPeriodStart, &payroll.PayPeriodEnd,
		&payroll.BaseSalary, &payroll.OvertimeHours, &payroll.OvertimePay, &payroll.HolidayHours,
		&payroll.HolidayPay, &payroll.Allowances, &payroll.Bonus, &payroll.GrossPay,
		&payroll.IncomeTax, &payroll.LocalTax, &payroll.NationalPension, &payroll.HealthInsurance,
		&payroll.EmploymentInsurance, &payroll.LongTermCare, &payroll.OtherDeductions,
		&payroll.TotalDeductions, &payroll.NetPay, &payroll.PayDate, &payroll.IsPaid,
		&payroll.CreatedAt, &payroll.UpdatedAt, &employeeName, &employeeNumber,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payroll record not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	payrollData := map[string]interface{}{
		"payroll":         payroll,
		"employee_name":   employeeName,
		"employee_number": employeeNumber,
	}

	c.JSON(http.StatusOK, payrollData)
}

func CreatePayrollRecord(c *gin.Context) {
	var req CreatePayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse pay period dates
	payPeriodStart, err := time.Parse("2006-01-02", req.PayPeriodStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pay period start date format (YYYY-MM-DD)"})
		return
	}

	payPeriodEnd, err := time.Parse("2006-01-02", req.PayPeriodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pay period end date format (YYYY-MM-DD)"})
		return
	}

	// Calculate payroll using the calculator
	calculator := PayrollCalculator{
		BaseSalary:      req.BaseSalary,
		OvertimeHours:   req.OvertimeHours,
		HolidayHours:    req.HolidayHours,
		Allowances:      req.Allowances,
		Bonus:           req.Bonus,
		OtherDeductions: req.OtherDeductions,
	}

	calculations := calculator.Calculate()

	// Insert payroll record
	result, err := database.DB.Exec(`
		INSERT INTO payroll_records (employee_id, pay_period_start, pay_period_end, 
		                            base_salary, overtime_hours, overtime_pay, holiday_hours, 
		                            holiday_pay, allowances, bonus, gross_pay, income_tax, 
		                            local_tax, national_pension, health_insurance, employment_insurance, 
		                            long_term_care, other_deductions, total_deductions, net_pay)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.EmployeeID, payPeriodStart, payPeriodEnd, req.BaseSalary, req.OvertimeHours,
		calculations["overtime_pay"], req.HolidayHours, calculations["holiday_pay"],
		req.Allowances, req.Bonus, calculations["gross_pay"], calculations["income_tax"],
		calculations["local_tax"], calculations["national_pension"], calculations["health_insurance"],
		calculations["employment_insurance"], calculations["long_term_care"], req.OtherDeductions,
		calculations["total_deductions"], calculations["net_pay"])

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payroll record"})
		return
	}

	payrollID, _ := result.LastInsertId()

	// Retrieve created payroll record with employee info
	var payroll models.PayrollRecord
	var employeeName, employeeNumber string

	err = database.DB.QueryRow(`
		SELECT p.id, p.employee_id, p.pay_period_start, p.pay_period_end, 
		       p.base_salary, p.overtime_hours, p.overtime_pay, p.holiday_hours, 
		       p.holiday_pay, p.allowances, p.bonus, p.gross_pay, p.income_tax, 
		       p.local_tax, p.national_pension, p.health_insurance, p.employment_insurance, 
		       p.long_term_care, p.other_deductions, p.total_deductions, p.net_pay, 
		       p.pay_date, p.is_paid, p.created_at, p.updated_at,
		       e.name as employee_name, e.employee_number
		FROM payroll_records p
		JOIN employees e ON p.employee_id = e.id
		WHERE p.id = ?
	`, payrollID).Scan(
		&payroll.ID, &payroll.EmployeeID, &payroll.PayPeriodStart, &payroll.PayPeriodEnd,
		&payroll.BaseSalary, &payroll.OvertimeHours, &payroll.OvertimePay, &payroll.HolidayHours,
		&payroll.HolidayPay, &payroll.Allowances, &payroll.Bonus, &payroll.GrossPay,
		&payroll.IncomeTax, &payroll.LocalTax, &payroll.NationalPension, &payroll.HealthInsurance,
		&payroll.EmploymentInsurance, &payroll.LongTermCare, &payroll.OtherDeductions,
		&payroll.TotalDeductions, &payroll.NetPay, &payroll.PayDate, &payroll.IsPaid,
		&payroll.CreatedAt, &payroll.UpdatedAt, &employeeName, &employeeNumber,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created payroll record"})
		return
	}

	payrollData := map[string]interface{}{
		"payroll":         payroll,
		"employee_name":   employeeName,
		"employee_number": employeeNumber,
	}

	c.JSON(http.StatusCreated, payrollData)
}

func UpdatePayrollRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return
	}

	var req CreatePayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse pay period dates
	payPeriodStart, err := time.Parse("2006-01-02", req.PayPeriodStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pay period start date format (YYYY-MM-DD)"})
		return
	}

	payPeriodEnd, err := time.Parse("2006-01-02", req.PayPeriodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pay period end date format (YYYY-MM-DD)"})
		return
	}

	// Recalculate payroll
	calculator := PayrollCalculator{
		BaseSalary:      req.BaseSalary,
		OvertimeHours:   req.OvertimeHours,
		HolidayHours:    req.HolidayHours,
		Allowances:      req.Allowances,
		Bonus:           req.Bonus,
		OtherDeductions: req.OtherDeductions,
	}

	calculations := calculator.Calculate()

	// Update payroll record
	_, err = database.DB.Exec(`
		UPDATE payroll_records SET pay_period_start = ?, pay_period_end = ?, 
		                          base_salary = ?, overtime_hours = ?, overtime_pay = ?, 
		                          holiday_hours = ?, holiday_pay = ?, allowances = ?, 
		                          bonus = ?, gross_pay = ?, income_tax = ?, local_tax = ?, 
		                          national_pension = ?, health_insurance = ?, employment_insurance = ?, 
		                          long_term_care = ?, other_deductions = ?, total_deductions = ?, 
		                          net_pay = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, payPeriodStart, payPeriodEnd, req.BaseSalary, req.OvertimeHours, calculations["overtime_pay"],
		req.HolidayHours, calculations["holiday_pay"], req.Allowances, req.Bonus,
		calculations["gross_pay"], calculations["income_tax"], calculations["local_tax"],
		calculations["national_pension"], calculations["health_insurance"], calculations["employment_insurance"],
		calculations["long_term_care"], req.OtherDeductions, calculations["total_deductions"],
		calculations["net_pay"], id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payroll record"})
		return
	}

	// Retrieve updated payroll record with employee info
	var payroll models.PayrollRecord
	var employeeName, employeeNumber string

	err = database.DB.QueryRow(`
		SELECT p.id, p.employee_id, p.pay_period_start, p.pay_period_end, 
		       p.base_salary, p.overtime_hours, p.overtime_pay, p.holiday_hours, 
		       p.holiday_pay, p.allowances, p.bonus, p.gross_pay, p.income_tax, 
		       p.local_tax, p.national_pension, p.health_insurance, p.employment_insurance, 
		       p.long_term_care, p.other_deductions, p.total_deductions, p.net_pay, 
		       p.pay_date, p.is_paid, p.created_at, p.updated_at,
		       e.name as employee_name, e.employee_number
		FROM payroll_records p
		JOIN employees e ON p.employee_id = e.id
		WHERE p.id = ?
	`, id).Scan(
		&payroll.ID, &payroll.EmployeeID, &payroll.PayPeriodStart, &payroll.PayPeriodEnd,
		&payroll.BaseSalary, &payroll.OvertimeHours, &payroll.OvertimePay, &payroll.HolidayHours,
		&payroll.HolidayPay, &payroll.Allowances, &payroll.Bonus, &payroll.GrossPay,
		&payroll.IncomeTax, &payroll.LocalTax, &payroll.NationalPension, &payroll.HealthInsurance,
		&payroll.EmploymentInsurance, &payroll.LongTermCare, &payroll.OtherDeductions,
		&payroll.TotalDeductions, &payroll.NetPay, &payroll.PayDate, &payroll.IsPaid,
		&payroll.CreatedAt, &payroll.UpdatedAt, &employeeName, &employeeNumber,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated payroll record"})
		return
	}

	payrollData := map[string]interface{}{
		"payroll":         payroll,
		"employee_name":   employeeName,
		"employee_number": employeeNumber,
	}

	c.JSON(http.StatusOK, payrollData)
}

func DeletePayrollRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return
	}

	// Delete payroll record
	_, err = database.DB.Exec("DELETE FROM payroll_records WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete payroll record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payroll record deleted successfully"})
}