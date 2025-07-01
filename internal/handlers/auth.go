package handlers

import (
	"database/sql"
	"labor-management-system/database"
	"labor-management-system/internal/middleware"
	"labor-management-system/internal/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// isPostgreSQL checks if we're using PostgreSQL
func isPostgreSQL() bool {
	return os.Getenv("DATABASE_URL") != ""
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	var query string
	if isPostgreSQL() {
		query = "SELECT id, username, password_hash, email, role, created_at, updated_at FROM users WHERE username = $1"
	} else {
		query = "SELECT id, username, password_hash, email, role, created_at, updated_at FROM users WHERE username = ?"
	}
	
	err := database.DB.QueryRow(query, req.Username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user,
	})
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "employee"
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert user into database
	var insertQuery string
	if isPostgreSQL() {
		insertQuery = "INSERT INTO users (username, password_hash, email, role) VALUES ($1, $2, $3, $4)"
	} else {
		insertQuery = "INSERT INTO users (username, password_hash, email, role) VALUES (?, ?, ?, ?)"
	}
	
	result, err := database.DB.Exec(insertQuery, req.Username, string(hashedPassword), req.Email, req.Role)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	userID, _ := result.LastInsertId()

	// Retrieve created user
	var user models.User
	var selectQuery string
	if isPostgreSQL() {
		selectQuery = "SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = $1"
	} else {
		selectQuery = "SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = ?"
	}
	
	err = database.DB.QueryRow(selectQuery, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, LoginResponse{
		Token: token,
		User:  user,
	})
}