package seeders

import (
	models "health/models/db"
	"health/services"

	"gorm.io/gorm"
)

func init() {
	services.RegisterSeeder(1, "user_seeder", func(database *gorm.DB) error {
		// Example: Create admin user
		admin := models.NewUser(
			"admin@example.com",
			"$2a$10$YourHashedPasswordHere", // Use bcrypt hash in production
			"Admin User",
			models.RoleAdmin,
		)

		// Check if user already exists
		var existingUser models.User
		result := database.Where("email = ?", admin.Email).First(&existingUser)
		if result.Error == nil {
			// User already exists, skip
			return nil
		}
		// If error is not "record not found", return it
		if result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}

		if err := database.Create(admin).Error; err != nil {
			return err
		}

		return nil
	})
}
