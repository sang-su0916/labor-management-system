package models

import (
	"database/sql"
	"time"
)

type Employee struct {
	ID             int            `json:"id" db:"id"`
	UserID         sql.NullInt64  `json:"user_id" db:"user_id"`
	EmployeeNumber string         `json:"employee_number" db:"employee_number"`
	Name           string         `json:"name" db:"name"`
	NameEn         sql.NullString `json:"name_en" db:"name_en"`
	Phone          sql.NullString `json:"phone" db:"phone"`
	Email          sql.NullString `json:"email" db:"email"`
	Address        sql.NullString `json:"address" db:"address"`
	BirthDate      sql.NullTime   `json:"birth_date" db:"birth_date"`
	HireDate       time.Time      `json:"hire_date" db:"hire_date"`
	Department     sql.NullString `json:"department" db:"department"`
	Position       sql.NullString `json:"position" db:"position"`
	EmploymentType string         `json:"employment_type" db:"employment_type"`
	Status         string         `json:"status" db:"status"`
	SalaryType     string         `json:"salary_type" db:"salary_type"`
	BaseSalary     sql.NullFloat64 `json:"base_salary" db:"base_salary"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

type EmploymentContract struct {
	ID             int            `json:"id" db:"id"`
	EmployeeID     int            `json:"employee_id" db:"employee_id"`
	ContractType   string         `json:"contract_type" db:"contract_type"`
	StartDate      time.Time      `json:"start_date" db:"start_date"`
	EndDate        sql.NullTime   `json:"end_date" db:"end_date"`
	Workplace      string         `json:"workplace" db:"workplace"`
	JobDescription sql.NullString `json:"job_description" db:"job_description"`
	WorkingHours   string         `json:"working_hours" db:"working_hours"`
	WorkDays       string         `json:"work_days" db:"work_days"`
	BaseSalary     float64        `json:"base_salary" db:"base_salary"`
	Allowances     sql.NullString `json:"allowances" db:"allowances"`
	Benefits       sql.NullString `json:"benefits" db:"benefits"`
	ContractTerms  sql.NullString `json:"contract_terms" db:"contract_terms"`
	SignedDate     sql.NullTime   `json:"signed_date" db:"signed_date"`
	IsActive       bool           `json:"is_active" db:"is_active"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

type AttendanceLog struct {
	ID            int            `json:"id" db:"id"`
	EmployeeID    int            `json:"employee_id" db:"employee_id"`
	WorkDate      time.Time      `json:"work_date" db:"work_date"`
	ClockIn       sql.NullString `json:"clock_in" db:"clock_in"`
	ClockOut      sql.NullString `json:"clock_out" db:"clock_out"`
	BreakStart    sql.NullString `json:"break_start" db:"break_start"`
	BreakEnd      sql.NullString `json:"break_end" db:"break_end"`
	TotalHours    sql.NullFloat64 `json:"total_hours" db:"total_hours"`
	OvertimeHours sql.NullFloat64 `json:"overtime_hours" db:"overtime_hours"`
	Status        string         `json:"status" db:"status"`
	Notes         sql.NullString `json:"notes" db:"notes"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
}

type PayrollRecord struct {
	ID                  int       `json:"id" db:"id"`
	EmployeeID          int       `json:"employee_id" db:"employee_id"`
	PayPeriodStart      time.Time `json:"pay_period_start" db:"pay_period_start"`
	PayPeriodEnd        time.Time `json:"pay_period_end" db:"pay_period_end"`
	BaseSalary          float64   `json:"base_salary" db:"base_salary"`
	OvertimeHours       float64   `json:"overtime_hours" db:"overtime_hours"`
	OvertimePay         float64   `json:"overtime_pay" db:"overtime_pay"`
	HolidayHours        float64   `json:"holiday_hours" db:"holiday_hours"`
	HolidayPay          float64   `json:"holiday_pay" db:"holiday_pay"`
	Allowances          float64   `json:"allowances" db:"allowances"`
	Bonus               float64   `json:"bonus" db:"bonus"`
	GrossPay            float64   `json:"gross_pay" db:"gross_pay"`
	IncomeTax           float64   `json:"income_tax" db:"income_tax"`
	LocalTax            float64   `json:"local_tax" db:"local_tax"`
	NationalPension     float64   `json:"national_pension" db:"national_pension"`
	HealthInsurance     float64   `json:"health_insurance" db:"health_insurance"`
	EmploymentInsurance float64   `json:"employment_insurance" db:"employment_insurance"`
	LongTermCare        float64   `json:"long_term_care" db:"long_term_care"`
	OtherDeductions     float64   `json:"other_deductions" db:"other_deductions"`
	TotalDeductions     float64   `json:"total_deductions" db:"total_deductions"`
	NetPay              float64   `json:"net_pay" db:"net_pay"`
	PayDate             sql.NullTime `json:"pay_date" db:"pay_date"`
	IsPaid              bool      `json:"is_paid" db:"is_paid"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type LeaveRequest struct {
	ID             int            `json:"id" db:"id"`
	EmployeeID     int            `json:"employee_id" db:"employee_id"`
	LeaveType      string         `json:"leave_type" db:"leave_type"`
	StartDate      time.Time      `json:"start_date" db:"start_date"`
	EndDate        time.Time      `json:"end_date" db:"end_date"`
	DaysRequested  float64        `json:"days_requested" db:"days_requested"`
	Reason         sql.NullString `json:"reason" db:"reason"`
	Status         string         `json:"status" db:"status"`
	ApprovedBy     sql.NullInt64  `json:"approved_by" db:"approved_by"`
	ApprovedAt     sql.NullTime   `json:"approved_at" db:"approved_at"`
	RejectionReason sql.NullString `json:"rejection_reason" db:"rejection_reason"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

type AnnualLeaveBalance struct {
	ID             int       `json:"id" db:"id"`
	EmployeeID     int       `json:"employee_id" db:"employee_id"`
	Year           int       `json:"year" db:"year"`
	TotalDays      float64   `json:"total_days" db:"total_days"`
	UsedDays       float64   `json:"used_days" db:"used_days"`
	RemainingDays  float64   `json:"remaining_days" db:"remaining_days"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}