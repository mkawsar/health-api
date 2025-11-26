package seeders

import (
	models "health/models/db"
	"health/services"

	"gorm.io/gorm"
)

func init() {
	services.RegisterSeeder(2, "doctor_seeder", func(database *gorm.DB) error {
		// Example: Create sample doctors
		doctors := []*models.Doctor{
			models.NewDoctor(
				"Dr. John Smith",
				"Cardiology",
				"+1234567890",
				"10 years",
				"New York",
				"LIC001",
				"9 AM - 5 PM",
				true,
				[]string{"Monday", "Wednesday", "Friday"},
				[]string{"09:00", "09:00", "09:00"},
				[]string{"17:00", "17:00", "17:00"},
			),
			models.NewDoctor(
				"Dr. Jane Doe",
				"Pediatrics",
				"+1234567891",
				"8 years",
				"Los Angeles",
				"LIC002",
				"10 AM - 6 PM",
				true,
				[]string{"Tuesday", "Thursday"},
				[]string{"10:00", "10:00"},
				[]string{"18:00", "18:00"},
			),
		}

		for _, doctor := range doctors {
			// Check if doctor already exists
			var existingDoctor models.Doctor
			result := database.Where("license = ?", doctor.License).First(&existingDoctor)
			if result.Error == nil {
				// Doctor already exists, skip
				continue
			}
			// If error is not "record not found", return it
			if result.Error != gorm.ErrRecordNotFound {
				return result.Error
			}

			if err := database.Create(doctor).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

