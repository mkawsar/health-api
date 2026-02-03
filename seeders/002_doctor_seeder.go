package seeders

import (
	models "health/models/db"
	"health/services"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	services.RegisterSeeder(2, "doctor_seeder", func() error {
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
			err := mgm.Coll(&models.Doctor{}).First(bson.M{"license": doctor.License}, &existingDoctor)
			if err == nil {
				// Doctor already exists, skip
				continue
			}

			// Create doctor
			if err := mgm.Coll(doctor).Create(doctor); err != nil {
				return err
			}
		}

		return nil
	})
}
