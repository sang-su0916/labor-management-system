-- 노무관리 시스템 데이터베이스 스키마

-- 사용자 (관리자, HR, 직원)
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    role VARCHAR(20) DEFAULT 'employee', -- admin, hr, employee
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 직원 정보
CREATE TABLE IF NOT EXISTS employees (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    employee_number VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    name_en VARCHAR(50),
    phone VARCHAR(20),
    email VARCHAR(100),
    address TEXT,
    birth_date DATE,
    hire_date DATE NOT NULL,
    department VARCHAR(50),
    position VARCHAR(50),
    employment_type VARCHAR(20) DEFAULT 'regular', -- regular, contract, part_time
    status VARCHAR(20) DEFAULT 'active', -- active, inactive, terminated
    salary_type VARCHAR(20) DEFAULT 'monthly', -- monthly, hourly, daily
    base_salary DECIMAL(10,2),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 근로계약서
CREATE TABLE IF NOT EXISTS employment_contracts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id INTEGER NOT NULL,
    contract_type VARCHAR(20) NOT NULL, -- regular, contract, part_time
    start_date DATE NOT NULL,
    end_date DATE,
    workplace TEXT NOT NULL,
    job_description TEXT,
    working_hours TEXT NOT NULL,
    work_days TEXT NOT NULL,
    base_salary DECIMAL(10,2) NOT NULL,
    allowances TEXT, -- JSON format
    benefits TEXT, -- JSON format
    contract_terms TEXT,
    signed_date DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id)
);

-- 급여 정보
CREATE TABLE IF NOT EXISTS payroll_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id INTEGER NOT NULL,
    pay_period_start DATE NOT NULL,
    pay_period_end DATE NOT NULL,
    base_salary DECIMAL(10,2) NOT NULL,
    overtime_hours DECIMAL(5,2) DEFAULT 0,
    overtime_pay DECIMAL(10,2) DEFAULT 0,
    holiday_hours DECIMAL(5,2) DEFAULT 0,
    holiday_pay DECIMAL(10,2) DEFAULT 0,
    allowances DECIMAL(10,2) DEFAULT 0, -- 각종 수당
    bonus DECIMAL(10,2) DEFAULT 0,
    gross_pay DECIMAL(10,2) NOT NULL, -- 총 지급액
    income_tax DECIMAL(10,2) DEFAULT 0, -- 소득세
    local_tax DECIMAL(10,2) DEFAULT 0, -- 지방소득세
    national_pension DECIMAL(10,2) DEFAULT 0, -- 국민연금
    health_insurance DECIMAL(10,2) DEFAULT 0, -- 건강보험
    employment_insurance DECIMAL(10,2) DEFAULT 0, -- 고용보험
    long_term_care DECIMAL(10,2) DEFAULT 0, -- 장기요양보험
    other_deductions DECIMAL(10,2) DEFAULT 0, -- 기타 공제
    total_deductions DECIMAL(10,2) NOT NULL, -- 총 공제액
    net_pay DECIMAL(10,2) NOT NULL, -- 실지급액
    pay_date DATE,
    is_paid BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id)
);

-- 근태 기록
CREATE TABLE IF NOT EXISTS attendance_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id INTEGER NOT NULL,
    work_date DATE NOT NULL,
    clock_in TIME,
    clock_out TIME,
    break_start TIME,
    break_end TIME,
    total_hours DECIMAL(4,2),
    overtime_hours DECIMAL(4,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'present', -- present, absent, late, early_leave, holiday
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    UNIQUE(employee_id, work_date)
);

-- 휴가 관리
CREATE TABLE IF NOT EXISTS leave_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id INTEGER NOT NULL,
    leave_type VARCHAR(20) NOT NULL, -- annual, sick, maternity, paternity, personal
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    days_requested DECIMAL(3,1) NOT NULL,
    reason TEXT,
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected, cancelled
    approved_by INTEGER,
    approved_at DATETIME,
    rejection_reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (approved_by) REFERENCES users(id)
);

-- 연차 잔여 현황
CREATE TABLE IF NOT EXISTS annual_leave_balance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id INTEGER NOT NULL,
    year INTEGER NOT NULL,
    total_days DECIMAL(3,1) NOT NULL, -- 총 연차 일수
    used_days DECIMAL(3,1) DEFAULT 0, -- 사용한 연차 일수
    remaining_days DECIMAL(3,1) NOT NULL, -- 잔여 연차 일수
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    UNIQUE(employee_id, year)
);

-- 문서 템플릿
CREATE TABLE IF NOT EXISTS document_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL, -- employment_contract, payslip, certificate, etc.
    content TEXT NOT NULL, -- HTML template
    variables TEXT, -- JSON format for template variables
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 생성된 문서 기록
CREATE TABLE IF NOT EXISTS generated_documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id INTEGER,
    template_id INTEGER NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    file_path VARCHAR(255),
    generated_by INTEGER NOT NULL,
    generated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (template_id) REFERENCES document_templates(id),
    FOREIGN KEY (generated_by) REFERENCES users(id)
);

-- 시스템 설정
CREATE TABLE IF NOT EXISTS system_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    setting_key VARCHAR(100) UNIQUE NOT NULL,
    setting_value TEXT,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 인덱스 생성
CREATE INDEX IF NOT EXISTS idx_employees_employee_number ON employees(employee_number);
CREATE INDEX IF NOT EXISTS idx_employees_department ON employees(department);
CREATE INDEX IF NOT EXISTS idx_employees_status ON employees(status);
CREATE INDEX IF NOT EXISTS idx_payroll_employee_period ON payroll_records(employee_id, pay_period_start, pay_period_end);
CREATE INDEX IF NOT EXISTS idx_attendance_employee_date ON attendance_logs(employee_id, work_date);
CREATE INDEX IF NOT EXISTS idx_leave_requests_employee ON leave_requests(employee_id);
CREATE INDEX IF NOT EXISTS idx_leave_requests_status ON leave_requests(status);

-- 기본 데이터 삽입
INSERT OR IGNORE INTO system_settings (setting_key, setting_value, description) VALUES
('company_name', '테스트 회사', '회사명'),
('company_address', '서울특별시 강남구 테헤란로 123', '회사 주소'),
('company_phone', '02-1234-5678', '회사 전화번호'),
('company_registration_number', '123-45-67890', '사업자등록번호'),
('min_wage', '9860', '최저임금 (시급)'),
('work_hours_per_day', '8', '1일 근무시간'),
('work_days_per_week', '5', '주 근무일수'),
('annual_leave_base', '15', '기본 연차 일수');

-- 관리자 계정 생성 (비밀번호: admin123!)
INSERT OR IGNORE INTO users (username, password_hash, email, role) VALUES
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin@company.com', 'admin');