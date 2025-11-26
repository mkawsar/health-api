package services

import (
	"context"
	"errors"
	db "health/models/db"
	"health/utils/requests"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user in the PostgreSQL database.
// The password is hashed using bcrypt.
// The user is created with the role "user".
// If the user cannot be created, an error is returned.
func CreateUser(name string, email string, password string) (*db.User, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("cannot generate hashed password")
	}

	user := db.NewUser(email, string(pass), name, db.RoleUser)
	err = DB.Create(user).Error
	if err != nil {
		return nil, errors.New("cannot create new user")
	}

	return user, nil
}

// FindUserById retrieves a user from the PostgreSQL database by the given ID.
// If the user does not exist, an error is returned.
func FindUserById(userId uint) (*db.User, error) {
	user := &db.User{}
	err := DB.First(user, userId).Error
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return user, nil
}

// FindUserByEmail retrieves a user from the PostgreSQL database by the given email address.
// If the user does not exist, an error is returned.
func FindUserByEmail(email string) (*db.User, error) {
	user := &db.User{}
	err := DB.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return user, nil
}

// CheckUserMail checks if a user with the given email address already exists in the PostgreSQL database.
// If such a user exists, an error is returned.
// If no such user exists, the function returns nil.
func CheckUserMail(email string) error {
	var count int64
	err := DB.Model(&db.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("email is already in use")
	}

	return nil
}

func GetUSers(ctx context.Context, page int, limit int, nameFilter string) ([]db.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var users []db.User
	query := DB.Model(&db.User{})
	
	if nameFilter != "" {
		query = query.Where("name ILIKE ?", "%"+nameFilter+"%") // Case-insensitive search
	}

	var totalUsers int64
	err := query.Count(&totalUsers).Error
	if err != nil {
		return nil, 0, errors.New("cannot count usersssss")
	}

	err = query.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, errors.New("cannot find users")
	}

	return users, totalUsers, nil
}

// GetUser retrieves a user from the PostgreSQL database by the given ID.
// If the user does not exist, an error is returned.
func GetUser(id uint) (*db.User, error) {
	user := &db.User{}
	err := DB.First(user, id).Error

	if err != nil {
		return nil, errors.New("cannot find user")
	}
	return user, nil
}

// UpdateUser updates the user's name in the PostgreSQL database with the given ID.
// The user's name is updated with the name provided in the UserRequest.
// If the user does not exist, an error is returned.
// If the user cannot be updated, an error is returned.
func UpdateUser(id uint, request *requests.UserRequest) error {
	user := &db.User{}
	err := DB.First(user, id).Error

	if err != nil {
		return errors.New("cannot find user")
	}
	user.Name = request.Name
	err = DB.Save(user).Error
	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

func DeleteUser(id uint) error {
	result := DB.Delete(&db.User{}, id)
	if result.Error != nil || result.RowsAffected == 0 {
		return errors.New("cannot delete user")
	}

	return nil
}