-- PostgreSQL용 스키마 (Railway 배포용)

-- 사용자 관리
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    role VARCHAR(20) DEFAULT 'employee',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 직원 정보
CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    employee_number VARCHAR(20) UNIQUE NOT NULL,
    department VARCHAR(50) NOT NULL,
    position VARCHAR(50) NOT NULL,
    hire_date DATE NOT NULL,
    salary DECIMAL(12,2) DEFAULT 0,
    phone VARCHAR(20),
    email VARCHAR(100),
    address TEXT,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 근로계약서
CREATE TABLE IF NOT EXISTS employment_contracts (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER NOT NULL REFERENCES employees(id),
    position VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    salary DECIMAL(12,2) NOT NULL,
    working_hours INTEGER DEFAULT 8,
    work_days VARCHAR(50) DEFAULT '월~금',
    contract_terms TEXT,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 급여 관리
CREATE TABLE IF NOT EXISTS payroll_records (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER NOT NULL REFERENCES employees(id),
    pay_period VARCHAR(7) NOT NULL,
    pay_period_start DATE,
    pay_period_end DATE,
    base_salary DECIMAL(12,2) NOT NULL,
    allowances DECIMAL(12,2) DEFAULT 0,
    bonus DECIMAL(12,2) DEFAULT 0,
    gross_pay DECIMAL(12,2) NOT NULL,
    income_tax DECIMAL(12,2) DEFAULT 0,
    local_tax DECIMAL(12,2) DEFAULT 0,
    national_pension DECIMAL(12,2) DEFAULT 0,
    health_insurance DECIMAL(12,2) DEFAULT 0,
    employment_insurance DECIMAL(12,2) DEFAULT 0,
    long_term_care DECIMAL(12,2) DEFAULT 0,
    other_deductions DECIMAL(12,2) DEFAULT 0,
    total_deductions DECIMAL(12,2) NOT NULL,
    net_pay DECIMAL(12,2) NOT NULL,
    pay_date DATE,
    is_paid BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 근태 기록
CREATE TABLE IF NOT EXISTS attendance_logs (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER NOT NULL REFERENCES employees(id),
    work_date DATE NOT NULL,
    clock_in TIME,
    clock_out TIME,
    break_start TIME,
    break_end TIME,
    total_hours DECIMAL(4,2),
    overtime_hours DECIMAL(4,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'present',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employee_id, work_date)
);

-- 휴가 관리
CREATE TABLE IF NOT EXISTS leave_requests (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER NOT NULL REFERENCES employees(id),
    leave_type VARCHAR(20) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    days_requested DECIMAL(3,1) NOT NULL,
    reason TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    approved_by INTEGER REFERENCES users(id),
    approved_at TIMESTAMP,
    rejection_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 연차 잔여 현황
CREATE TABLE IF NOT EXISTS annual_leave_balance (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER NOT NULL REFERENCES employees(id),
    year INTEGER NOT NULL,
    total_days DECIMAL(3,1) NOT NULL,
    used_days DECIMAL(3,1) DEFAULT 0,
    remaining_days DECIMAL(3,1) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employee_id, year)
);

-- 문서 템플릿
CREATE TABLE IF NOT EXISTS document_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    variables TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 생성된 문서 기록
CREATE TABLE IF NOT EXISTS generated_documents (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER REFERENCES employees(id),
    template_id INTEGER NOT NULL REFERENCES document_templates(id),
    document_type VARCHAR(50) NOT NULL,
    file_path VARCHAR(255),
    generated_by INTEGER NOT NULL REFERENCES users(id),
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 시스템 설정
CREATE TABLE IF NOT EXISTS system_settings (
    id SERIAL PRIMARY KEY,
    setting_key VARCHAR(100) UNIQUE NOT NULL,
    setting_value TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
INSERT INTO system_settings (setting_key, setting_value, description) VALUES
('company_name', '테스트 회사', '회사명'),
('company_address', '서울특별시 강남구 테헤란로 123', '회사 주소'),
('company_phone', '02-1234-5678', '회사 전화번호'),
('company_registration_number', '123-45-67890', '사업자등록번호'),
('min_wage', '9860', '최저임금 (시급)'),
('work_hours_per_day', '8', '1일 근무시간'),
('work_days_per_week', '5', '주 근무일수'),
('annual_leave_base', '15', '기본 연차 일수')
ON CONFLICT (setting_key) DO NOTHING;

-- 관리자 계정 생성 (비밀번호: admin123)
INSERT INTO users (username, password_hash, email, role) VALUES
('admin', '$2a$10$BGuuHyAsIfgXDObMqhNUwOnfY4oK56B50BVx1NoZWL0y9kRmsdYji', 'admin@company.com', 'admin')
ON CONFLICT (username) DO NOTHING;