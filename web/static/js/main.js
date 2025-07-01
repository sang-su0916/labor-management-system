// Labor Management System - Complete Implementation

// Global variables
let authToken = null;
let currentUser = null;
let currentEmployees = [];

// API Base URL
const API_BASE = '/api';

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
    setupEventListeners();
});

function initializeApp() {
    console.log('Initializing Labor Management System...');
    
    // Check if user is already logged in
    authToken = localStorage.getItem('authToken');
    currentUser = JSON.parse(localStorage.getItem('currentUser') || 'null');
    
    if (authToken && currentUser) {
        showDashboard();
    } else {
        showLogin();
    }
    
    // Setup login form
    setupLoginForm();
}

function setupEventListeners() {
    // Date change listeners for leave request
    document.addEventListener('change', function(e) {
        if (e.target.id === 'startDate' || e.target.id === 'endDate') {
            calculateLeaveDays();
        }
        if (e.target.id === 'payrollEmployee') {
            loadEmployeeSalary();
        }
    });
}

function setupLoginForm() {
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }
}

async function handleLogin(event) {
    event.preventDefault();
    
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    
    if (!username || !password) {
        showAlert('사용자명과 비밀번호를 입력해주세요.', 'warning');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                password: password
            })
        });
        
        const data = await response.json();
        
        if (response.ok && data.token) {
            authToken = data.token;
            currentUser = data.user;
            
            // Store in localStorage
            localStorage.setItem('authToken', authToken);
            localStorage.setItem('currentUser', JSON.stringify(currentUser));
            
            showAlert('로그인이 성공했습니다!', 'success');
            showDashboard();
        } else {
            showAlert(data.error || '로그인에 실패했습니다.', 'danger');
        }
    } catch (error) {
        console.error('Login error:', error);
        showAlert('서버 연결에 실패했습니다.', 'danger');
    }
}

function showLogin() {
    document.getElementById('loginSection').classList.remove('d-none');
    document.getElementById('dashboardSection').classList.add('d-none');
}

function showDashboard() {
    document.getElementById('loginSection').classList.add('d-none');
    document.getElementById('dashboardSection').classList.remove('d-none');
    document.getElementById('dashboardSection').classList.add('fade-in');
}

function logout() {
    authToken = null;
    currentUser = null;
    localStorage.removeItem('authToken');
    localStorage.removeItem('currentUser');
    
    showAlert('로그아웃되었습니다.', 'info');
    showLogin();
    
    // Clear content area
    document.getElementById('contentArea').innerHTML = '';
}

// Alert function
function showAlert(message, type = 'info', duration = 5000) {
    const alertContainer = document.getElementById('alertContainer');
    const alertId = 'alert-' + Date.now();
    
    const alertHTML = `
        <div id="${alertId}" class="alert alert-${type} alert-dismissible fade show" role="alert">
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        </div>
    `;
    
    alertContainer.insertAdjacentHTML('beforeend', alertHTML);
    
    // Auto-dismiss after duration
    setTimeout(() => {
        const alertElement = document.getElementById(alertId);
        if (alertElement) {
            const bsAlert = new bootstrap.Alert(alertElement);
            bsAlert.close();
        }
    }, duration);
}

// API helper function
async function apiCall(endpoint, options = {}) {
    const config = {
        headers: {
            'Content-Type': 'application/json',
            ...options.headers
        },
        ...options
    };
    
    if (authToken) {
        config.headers['Authorization'] = `Bearer ${authToken}`;
    }
    
    try {
        const response = await fetch(`${API_BASE}${endpoint}`, config);
        const data = await response.json();
        
        if (!response.ok) {
            if (response.status === 401) {
                // Token expired, logout
                logout();
                throw new Error('인증이 만료되었습니다. 다시 로그인해주세요.');
            }
            throw new Error(data.error || '요청 처리 중 오류가 발생했습니다.');
        }
        
        return data;
    } catch (error) {
        console.error('API call error:', error);
        throw error;
    }
}

// Load functions for different modules
async function loadEmployees() {
    try {
        const data = await apiCall('/employees');
        currentEmployees = data.employees || [];
        displayEmployees(currentEmployees);
        updateEmployeeSelects();
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

async function loadAttendance() {
    try {
        const data = await apiCall('/attendance');
        displayAttendance(data.attendance || []);
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

async function loadLeaves() {
    try {
        const data = await apiCall('/leaves');
        displayLeaves(data.leaves || []);
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

async function loadPayroll() {
    try {
        const data = await apiCall('/payroll');
        displayPayroll(data.payroll || []);
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

// Update employee select options
function updateEmployeeSelects() {
    const payrollSelect = document.getElementById('payrollEmployee');
    const leaveSelect = document.getElementById('leaveEmployee');
    
    if (payrollSelect) {
        payrollSelect.innerHTML = '<option value="">직원 선택</option>';
        currentEmployees.forEach(emp => {
            payrollSelect.innerHTML += `<option value="${emp.id}">${emp.name} (${emp.employee_number})</option>`;
        });
    }
    
    if (leaveSelect) {
        leaveSelect.innerHTML = '<option value="">직원 선택</option>';
        currentEmployees.forEach(emp => {
            leaveSelect.innerHTML += `<option value="${emp.id}">${emp.name} (${emp.employee_number})</option>`;
        });
    }
}

// Display functions
function displayEmployees(employees) {
    const contentArea = document.getElementById('contentArea');
    
    let html = `
        <div class="d-flex justify-content-between align-items-center mb-3">
            <h4>직원 목록</h4>
            <button class="btn btn-primary" onclick="showAddEmployeeModal()">
                <i class="fas fa-plus"></i> 직원 추가
            </button>
        </div>
        <div class="table-responsive">
            <table class="table table-striped">
                <thead>
                    <tr>
                        <th>사번</th>
                        <th>이름</th>
                        <th>부서</th>
                        <th>직급</th>
                        <th>입사일</th>
                        <th>상태</th>
                        <th>액션</th>
                    </tr>
                </thead>
                <tbody>
    `;
    
    employees.forEach(employee => {
        html += `
            <tr>
                <td>${employee.employee_number}</td>
                <td>${employee.name}</td>
                <td>${employee.department}</td>
                <td>${employee.position}</td>
                <td>${formatDate(employee.hire_date)}</td>
                <td><span class="status-badge status-${employee.status}">${employee.status}</span></td>
                <td>
                    <button class="btn btn-sm btn-outline-primary" onclick="viewEmployee(${employee.id})">보기</button>
                    <button class="btn btn-sm btn-outline-warning" onclick="editEmployee(${employee.id})">수정</button>
                    <button class="btn btn-sm btn-outline-danger" onclick="deleteEmployee(${employee.id})">삭제</button>
                </td>
            </tr>
        `;
    });
    
    html += `
                </tbody>
            </table>
        </div>
    `;
    
    contentArea.innerHTML = html;
}

function displayAttendance(attendance) {
    const contentArea = document.getElementById('contentArea');
    
    let html = `
        <div class="d-flex justify-content-between align-items-center mb-3">
            <h4>근태 현황</h4>
            <div>
                <button class="btn btn-success" onclick="clockIn()">출근</button>
                <button class="btn btn-warning" onclick="clockOut()">퇴근</button>
            </div>
        </div>
        <div class="table-responsive">
            <table class="table table-striped">
                <thead>
                    <tr>
                        <th>직원명</th>
                        <th>날짜</th>
                        <th>출근시간</th>
                        <th>퇴근시간</th>
                        <th>총 근무시간</th>
                        <th>상태</th>
                    </tr>
                </thead>
                <tbody>
    `;
    
    attendance.forEach(record => {
        html += `
            <tr>
                <td>${record.employee_name}</td>
                <td>${formatDate(record.work_date)}</td>
                <td>${record.clock_in || '-'}</td>
                <td>${record.clock_out || '-'}</td>
                <td>${record.total_hours || '-'}시간</td>
                <td><span class="status-badge status-${record.status}">${record.status}</span></td>
            </tr>
        `;
    });
    
    html += `
                </tbody>
            </table>
        </div>
    `;
    
    contentArea.innerHTML = html;
}

function displayLeaves(leaves) {
    const contentArea = document.getElementById('contentArea');
    
    let html = `
        <div class="d-flex justify-content-between align-items-center mb-3">
            <h4>휴가 현황</h4>
            <button class="btn btn-primary" onclick="showAddLeaveModal()">
                <i class="fas fa-plus"></i> 휴가 신청
            </button>
        </div>
        <div class="table-responsive">
            <table class="table table-striped">
                <thead>
                    <tr>
                        <th>직원명</th>
                        <th>휴가 유형</th>
                        <th>시작일</th>
                        <th>종료일</th>
                        <th>일수</th>
                        <th>상태</th>
                        <th>액션</th>
                    </tr>
                </thead>
                <tbody>
    `;
    
    leaves.forEach(leave => {
        html += `
            <tr>
                <td>${leave.employee_name}</td>
                <td>${getLeaveTypeLabel(leave.leave_type)}</td>
                <td>${formatDate(leave.start_date)}</td>
                <td>${formatDate(leave.end_date)}</td>
                <td>${leave.days_requested}일</td>
                <td><span class="status-badge status-${leave.status}">${getStatusLabel(leave.status)}</span></td>
                <td>
                    ${leave.status === 'pending' ? `
                        <button class="btn btn-sm btn-success" onclick="approveLeave(${leave.id})">승인</button>
                        <button class="btn btn-sm btn-danger" onclick="rejectLeave(${leave.id})">거절</button>
                    ` : ''}
                </td>
            </tr>
        `;
    });
    
    html += `
                </tbody>
            </table>
        </div>
    `;
    
    contentArea.innerHTML = html;
}

function displayPayroll(payroll) {
    const contentArea = document.getElementById('contentArea');
    
    let html = `
        <div class="d-flex justify-content-between align-items-center mb-3">
            <h4>급여 현황</h4>
            <button class="btn btn-primary" onclick="showPayrollModal()">
                <i class="fas fa-plus"></i> 급여 등록
            </button>
        </div>
        <div class="table-responsive">
            <table class="table table-striped">
                <thead>
                    <tr>
                        <th>직원명</th>
                        <th>급여기간</th>
                        <th>기본급</th>
                        <th>총 지급액</th>
                        <th>총 공제액</th>
                        <th>실지급액</th>
                        <th>지급상태</th>
                        <th>액션</th>
                    </tr>
                </thead>
                <tbody>
    `;
    
    payroll.forEach(record => {
        html += `
            <tr>
                <td>${record.employee_name}</td>
                <td>${record.pay_period}</td>
                <td>${formatCurrency(record.base_salary)}</td>
                <td>${formatCurrency(record.gross_pay)}</td>
                <td>${formatCurrency(record.total_deductions)}</td>
                <td><strong>${formatCurrency(record.net_pay)}</strong></td>
                <td><span class="status-badge ${record.is_paid ? 'status-approved' : 'status-pending'}">${record.is_paid ? '지급완료' : '지급대기'}</span></td>
                <td>
                    <button class="btn btn-sm btn-outline-primary" onclick="generatePayslip(${record.employee_id})">명세서</button>
                    <button class="btn btn-sm btn-outline-warning" onclick="editPayroll(${record.id})">수정</button>
                    <button class="btn btn-sm btn-outline-danger" onclick="deletePayroll(${record.id})">삭제</button>
                </td>
            </tr>
        `;
    });
    
    html += `
                </tbody>
            </table>
        </div>
    `;
    
    contentArea.innerHTML = html;
}

// Modal functions
function showAddEmployeeModal() {
    if (currentEmployees.length === 0) {
        loadEmployees().then(() => {
            updateEmployeeSelects();
        });
    }
    
    document.getElementById('employeeModalTitle').textContent = '직원 추가';
    document.getElementById('employeeForm').reset();
    document.getElementById('employeeId').value = '';
    
    const modal = new bootstrap.Modal(document.getElementById('employeeModal'));
    modal.show();
}

function showPayrollModal() {
    if (currentEmployees.length === 0) {
        loadEmployees().then(() => {
            updateEmployeeSelects();
        });
    } else {
        updateEmployeeSelects();
    }
    
    document.getElementById('payrollForm').reset();
    document.getElementById('payrollCalculation').classList.add('d-none');
    
    // Set current month
    const now = new Date();
    document.getElementById('payPeriod').value = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
    
    const modal = new bootstrap.Modal(document.getElementById('payrollModal'));
    modal.show();
}

function showAddLeaveModal() {
    if (currentEmployees.length === 0) {
        loadEmployees().then(() => {
            updateEmployeeSelects();
        });
    } else {
        updateEmployeeSelects();
    }
    
    document.getElementById('leaveForm').reset();
    
    const modal = new bootstrap.Modal(document.getElementById('leaveModal'));
    modal.show();
}

// Employee functions
async function saveEmployee() {
    const form = document.getElementById('employeeForm');
    const employeeId = document.getElementById('employeeId').value;
    
    const employeeData = {
        name: document.getElementById('employeeName').value,
        employee_number: document.getElementById('employeeNumber').value,
        department: document.getElementById('department').value,
        position: document.getElementById('position').value,
        hire_date: document.getElementById('hireDate').value,
        salary: parseFloat(document.getElementById('salary').value) || 0,
        phone: document.getElementById('phone').value,
        email: document.getElementById('email').value,
        address: document.getElementById('address').value
    };
    
    if (!employeeData.name || !employeeData.employee_number || !employeeData.department || 
        !employeeData.position || !employeeData.hire_date) {
        showAlert('필수 항목을 모두 입력해주세요.', 'warning');
        return;
    }
    
    try {
        const method = employeeId ? 'PUT' : 'POST';
        const endpoint = employeeId ? `/employees/${employeeId}` : '/employees';
        
        await apiCall(endpoint, {
            method: method,
            body: JSON.stringify(employeeData)
        });
        
        showAlert(`직원이 성공적으로 ${employeeId ? '수정' : '추가'}되었습니다.`, 'success');
        
        // Close modal and refresh list
        bootstrap.Modal.getInstance(document.getElementById('employeeModal')).hide();
        loadEmployees();
        
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

async function editEmployee(id) {
    try {
        const data = await apiCall(`/employees/${id}`);
        const employee = data.employee;
        
        document.getElementById('employeeModalTitle').textContent = '직원 수정';
        document.getElementById('employeeId').value = employee.id;
        document.getElementById('employeeName').value = employee.name;
        document.getElementById('employeeNumber').value = employee.employee_number;
        document.getElementById('department').value = employee.department;
        document.getElementById('position').value = employee.position;
        document.getElementById('hireDate').value = employee.hire_date;
        document.getElementById('salary').value = employee.salary || '';
        document.getElementById('phone').value = employee.phone || '';
        document.getElementById('email').value = employee.email || '';
        document.getElementById('address').value = employee.address || '';
        
        const modal = new bootstrap.Modal(document.getElementById('employeeModal'));
        modal.show();
        
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

async function deleteEmployee(id) {
    if (!confirm('정말로 이 직원을 삭제하시겠습니까?')) {
        return;
    }
    
    try {
        await apiCall(`/employees/${id}`, { method: 'DELETE' });
        showAlert('직원이 성공적으로 삭제되었습니다.', 'success');
        loadEmployees();
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

// Payroll functions
function calculatePayroll() {
    const baseSalary = parseFloat(document.getElementById('baseSalary').value) || 0;
    const allowances = parseFloat(document.getElementById('allowances').value) || 0;
    const bonus = parseFloat(document.getElementById('bonus').value) || 0;
    
    if (baseSalary === 0) {
        showAlert('기본급을 입력해주세요.', 'warning');
        return;
    }
    
    const grossPay = baseSalary + allowances + bonus;
    
    // Korean tax calculations (simplified)
    const incomeTax = grossPay * 0.03;
    const localTax = incomeTax * 0.1;
    const nationalPension = grossPay * 0.045;
    const healthInsurance = grossPay * 0.03335;
    const longTermCare = healthInsurance * 0.1281;
    const employmentInsurance = grossPay * 0.008;
    
    const totalDeductions = incomeTax + localTax + nationalPension + healthInsurance + longTermCare + employmentInsurance;
    const netPay = grossPay - totalDeductions;
    
    // Update fields
    document.getElementById('grossPay').value = Math.round(grossPay);
    document.getElementById('incomeTax').value = Math.round(incomeTax);
    document.getElementById('nationalPension').value = Math.round(nationalPension);
    document.getElementById('healthInsurance').value = Math.round(healthInsurance);
    document.getElementById('employmentInsurance').value = Math.round(employmentInsurance);
    document.getElementById('totalDeductions').value = Math.round(totalDeductions);
    document.getElementById('netPay').value = Math.round(netPay);
    
    document.getElementById('payrollCalculation').classList.remove('d-none');
}

async function savePayroll() {
    const payrollData = {
        employee_id: parseInt(document.getElementById('payrollEmployee').value),
        pay_period: document.getElementById('payPeriod').value,
        base_salary: parseFloat(document.getElementById('baseSalary').value),
        allowances: parseFloat(document.getElementById('allowances').value) || 0,
        bonus: parseFloat(document.getElementById('bonus').value) || 0,
        gross_pay: parseFloat(document.getElementById('grossPay').value),
        income_tax: parseFloat(document.getElementById('incomeTax').value),
        national_pension: parseFloat(document.getElementById('nationalPension').value),
        health_insurance: parseFloat(document.getElementById('healthInsurance').value),
        employment_insurance: parseFloat(document.getElementById('employmentInsurance').value),
        total_deductions: parseFloat(document.getElementById('totalDeductions').value),
        net_pay: parseFloat(document.getElementById('netPay').value)
    };
    
    if (!payrollData.employee_id || !payrollData.pay_period || !payrollData.base_salary) {
        showAlert('필수 항목을 모두 입력하고 급여를 계산해주세요.', 'warning');
        return;
    }
    
    try {
        await apiCall('/payroll', {
            method: 'POST',
            body: JSON.stringify(payrollData)
        });
        
        showAlert('급여가 성공적으로 등록되었습니다.', 'success');
        
        // Close modal and refresh list
        bootstrap.Modal.getInstance(document.getElementById('payrollModal')).hide();
        loadPayroll();
        
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

async function loadEmployeeSalary() {
    const employeeId = document.getElementById('payrollEmployee').value;
    if (!employeeId) return;
    
    try {
        const data = await apiCall(`/employees/${employeeId}`);
        const employee = data.employee;
        
        if (employee.salary) {
            document.getElementById('baseSalary').value = employee.salary;
        }
    } catch (error) {
        console.error('Failed to load employee salary:', error);
    }
}

// Leave functions
function calculateLeaveDays() {
    const startDate = document.getElementById('startDate').value;
    const endDate = document.getElementById('endDate').value;
    
    if (startDate && endDate) {
        const start = new Date(startDate);
        const end = new Date(endDate);
        const diffTime = Math.abs(end - start);
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24)) + 1;
        
        document.getElementById('daysRequested').value = diffDays;
    }
}

async function saveLeaveRequest() {
    const leaveData = {
        employee_id: parseInt(document.getElementById('leaveEmployee').value),
        leave_type: document.getElementById('leaveType').value,
        start_date: document.getElementById('startDate').value,
        end_date: document.getElementById('endDate').value,
        days_requested: parseFloat(document.getElementById('daysRequested').value),
        reason: document.getElementById('leaveReason').value
    };
    
    if (!leaveData.employee_id || !leaveData.leave_type || !leaveData.start_date || 
        !leaveData.end_date || !leaveData.days_requested) {
        showAlert('필수 항목을 모두 입력해주세요.', 'warning');
        return;
    }
    
    try {
        await apiCall('/leaves', {
            method: 'POST',
            body: JSON.stringify(leaveData)
        });
        
        showAlert('휴가 신청이 성공적으로 등록되었습니다.', 'success');
        
        // Close modal and refresh list
        bootstrap.Modal.getInstance(document.getElementById('leaveModal')).hide();
        loadLeaves();
        
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

async function approveLeave(id) {
    if (!confirm('이 휴가 신청을 승인하시겠습니까?')) {
        return;
    }
    
    try {
        await apiCall(`/leaves/${id}/approve`, {
            method: 'PUT',
            body: JSON.stringify({ approved_by: currentUser.id })
        });
        
        showAlert('휴가 신청이 승인되었습니다.', 'success');
        loadLeaves();
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

async function rejectLeave(id) {
    const reason = prompt('거절 사유를 입력해주세요:');
    if (!reason) return;
    
    try {
        await apiCall(`/leaves/${id}/reject`, {
            method: 'PUT',
            body: JSON.stringify({ 
                approved_by: currentUser.id,
                rejection_reason: reason 
            })
        });
        
        showAlert('휴가 신청이 거절되었습니다.', 'success');
        loadLeaves();
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

// Attendance functions
async function clockIn() {
    try {
        await apiCall('/attendance/clock-in', {
            method: 'POST',
            body: JSON.stringify({ employee_id: currentUser.id })
        });
        
        showAlert('출근이 기록되었습니다.', 'success');
        loadAttendance();
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

async function clockOut() {
    try {
        await apiCall('/attendance/clock-out', {
            method: 'POST',
            body: JSON.stringify({ employee_id: currentUser.id })
        });
        
        showAlert('퇴근이 기록되었습니다.', 'success');
        loadAttendance();
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

// Document functions
async function generatePayslip(employeeId) {
    try {
        const response = await fetch(`${API_BASE}/documents/generate/payslip?employee_id=${employeeId}`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${authToken}`,
                'Content-Type': 'application/json'
            }
        });
        
        if (response.ok) {
            const data = await response.json();
            showAlert('급여명세서가 생성되었습니다.', 'success');
        } else {
            const data = await response.json();
            throw new Error(data.error);
        }
    } catch (error) {
        showAlert(error.message, 'danger');
    }
}

// Utility functions
function formatCurrency(amount) {
    return new Intl.NumberFormat('ko-KR', {
        style: 'currency',
        currency: 'KRW'
    }).format(amount);
}

function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString('ko-KR');
}

function getLeaveTypeLabel(type) {
    const types = {
        'annual': '연차',
        'sick': '병가',
        'personal': '개인사유',
        'maternity': '출산휴가',
        'paternity': '육아휴직'
    };
    return types[type] || type;
}

function getStatusLabel(status) {
    const statuses = {
        'pending': '대기중',
        'approved': '승인됨',
        'rejected': '거절됨',
        'active': '활성',
        'inactive': '비활성'
    };
    return statuses[status] || status;
}

function viewEmployee(id) {
    editEmployee(id);
}