package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Email        string    `json:"email" db:"email"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type DocumentTemplate struct {
	ID        int            `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	Type      string         `json:"type" db:"type"`
	Content   string         `json:"content" db:"content"`
	Variables sql.NullString `json:"variables" db:"variables"`
	IsActive  bool           `json:"is_active" db:"is_active"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

type GeneratedDocument struct {
	ID           int            `json:"id" db:"id"`
	EmployeeID   sql.NullInt64  `json:"employee_id" db:"employee_id"`
	TemplateID   int            `json:"template_id" db:"template_id"`
	DocumentType string         `json:"document_type" db:"document_type"`
	FilePath     sql.NullString `json:"file_path" db:"file_path"`
	GeneratedBy  int            `json:"generated_by" db:"generated_by"`
	GeneratedAt  time.Time      `json:"generated_at" db:"generated_at"`
}

type SystemSetting struct {
	ID           int            `json:"id" db:"id"`
	SettingKey   string         `json:"setting_key" db:"setting_key"`
	SettingValue sql.NullString `json:"setting_value" db:"setting_value"`
	Description  sql.NullString `json:"description" db:"description"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
}