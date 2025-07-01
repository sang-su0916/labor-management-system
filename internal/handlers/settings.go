package handlers

import (
	"labor-management-system/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SystemSetting struct {
	ID           int    `json:"id"`
	SettingKey   string `json:"setting_key"`
	SettingValue string `json:"setting_value"`
	Description  string `json:"description"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type UpdateSettingsRequest struct {
	Settings map[string]string `json:"settings"`
}

func GetSystemSettings(c *gin.Context) {
	db := database.GetDB()
	
	rows, err := db.Query(`
		SELECT id, setting_key, COALESCE(setting_value, ''), 
		       COALESCE(description, ''), created_at, updated_at
		FROM system_settings
		ORDER BY setting_key
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	}
	defer rows.Close()

	settings := make(map[string]SystemSetting)
	for rows.Next() {
		var setting SystemSetting
		err := rows.Scan(&setting.ID, &setting.SettingKey, &setting.SettingValue,
			&setting.Description, &setting.CreatedAt, &setting.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan setting"})
			return
		}
		settings[setting.SettingKey] = setting
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

func UpdateSystemSettings(c *gin.Context) {
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	db := database.GetDB()
	
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Update each setting
	for key, value := range req.Settings {
		_, err := tx.Exec(`
			INSERT INTO system_settings (setting_key, setting_value, updated_at)
			VALUES (?, ?, CURRENT_TIMESTAMP)
			ON CONFLICT(setting_key) DO UPDATE SET
				setting_value = excluded.setting_value,
				updated_at = CURRENT_TIMESTAMP
		`, key, value)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update setting: " + key,
			})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit changes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

// Helper function to get a specific setting value
func GetSettingValue(key string) (string, error) {
	db := database.GetDB()
	
	var value string
	err := db.QueryRow(`
		SELECT COALESCE(setting_value, '') 
		FROM system_settings 
		WHERE setting_key = ?
	`, key).Scan(&value)
	
	return value, err
}

// Helper function to set a specific setting value
func SetSettingValue(key, value, description string) error {
	db := database.GetDB()
	
	_, err := db.Exec(`
		INSERT INTO system_settings (setting_key, setting_value, description, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(setting_key) DO UPDATE SET
			setting_value = excluded.setting_value,
			description = excluded.description,
			updated_at = CURRENT_TIMESTAMP
	`, key, value, description)
	
	return err
}