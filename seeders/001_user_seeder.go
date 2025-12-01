package seeders

import (
	"errors"
	models "health/models/db"
	"health/services"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	services.RegisterSeeder(1, "user_seeder", func() error {
		// Example: Create admin user
		admin := models.NewUser(
			"admin@example.com",
			"$2a$10$YourHashedPasswordHere", // Use bcrypt hash in production
			"Admin User",
			models.RoleAdmin,
		)

		// Check if user already exists
		var existingUser models.User
		err := mgm.Coll(&models.User{}).First(bson.M{"email": admin.Email}, &existingUser)
		if err == nil {
			// User already exists, skip
			return nil
		}

		// Create user
		if err := mgm.Coll(admin).Create(admin); err != nil {
			return errors.New("cannot create admin user")
		}

		return nil
	})
}
