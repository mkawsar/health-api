package services

import (
	"context"
	"errors"
	db "health/models/db"
	"health/utils/requests"

	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user in the MongoDB database.
// The password is hashed using bcrypt.
// The user is created with the role "user".
// If the user cannot be created, an error is returned.
func CreateUser(name string, email string, password string) (*db.User, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("cannot generate hashed password")
	}

	user := db.NewUser(email, string(pass), name, db.RoleUser)
	err = mgm.Coll(user).Create(user)
	if err != nil {
		return nil, errors.New("cannot create new user")
	}

	return user, nil
}

// FindUserById retrieves a user from the MongoDB database by the given ObjectID.
// If the user does not exist, an error is returned.
func FindUserById(userId primitive.ObjectID) (*db.User, error) {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(userId, user)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return user, nil
}

// FindUserByEmail retrieves a user from the MongoDB database by the given email address.
// If the user does not exist, an error is returned.
func FindUserByEmail(email string) (*db.User, error) {
	user := &db.User{}
	err := mgm.Coll(user).First(bson.M{"email": email}, user)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return user, nil
}

// CheckUserMail checks if a user with the given email address already exists in the MongoDB database.
// If such a user exists, an error is returned.
// If no such user exists, the function returns nil.
func CheckUserMail(email string) error {
	user := &db.User{}
	userCollection := mgm.Coll(user)
	err := userCollection.First(bson.M{"email": email}, user)
	if err == nil {
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
	skip := (page - 1) * limit
	filter := bson.M{}
	if nameFilter != "" {
		filter["name"] = bson.M{"$regex": nameFilter, "$options": "i"} // Case-insensitive search
	}
	var users []db.User
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(skip))
	err := mgm.Coll(&db.User{}).SimpleFind(&users, filter, opts)

	if err != nil {
		return nil, 0, errors.New("cannot find users")
	}
	totalUsers, _ := mgm.Coll(&db.User{}).CountDocuments(ctx, filter)
	return users, totalUsers, nil
}

// GetUser retrieves a user from the MongoDB database by the given ObjectID.
// If the user does not exist, an error is returned.
func GetUser(id primitive.ObjectID) (*db.User, error) {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(id, user)

	if err != nil {
		return nil, errors.New("cannot find user")
	}
	return user, nil
}

// UpdateUser updates the user's name in the MongoDB database with the given ObjectID.
// The user's name is updated with the name provided in the UserRequest.
// If the user does not exist, an error is returned.
// If the user cannot be updated, an error is returned.
func UpdateUser(id primitive.ObjectID, request *requests.UserRequest) error {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(id, user)

	if err != nil {
		return errors.New("cannot find user")
	}
	user.Name = request.Name
	err = mgm.Coll(user).Update(user)
	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

func DeleteUser(id primitive.ObjectID) error {
	result, err := mgm.Coll(&db.User{}).DeleteOne(mgm.Ctx(), bson.M{field.ID: id})
	if err != nil || result.DeletedCount <= 0 {
		return errors.New("cannot delete note")
	}

	return nil
}