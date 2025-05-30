package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/kamva/mgm/v3"
)

const (
	RoleUser = "user"
	RoleAdmin = "admin"
	RoleDoctor = "doctor"
	RolePharmacy = "pharmacy"
	RoleLab = "lab"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string `json:"email" bson:"email"`
	Password         string `json:"-" bson:"password"`
	Name             string `json:"name" bson:"name"`
	Role             string `json:"role" bson:"role"`
	MailVarified     bool   `json:"mail_verified" bson:"mail_verified"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Type  string `json:"type"`
}

func NewUser(email string, password string, name string, role string) *User {
	return &User{
		Email:        email,
		Password:     password,
		Name:         name,
		Role:         role,
		MailVarified: false,
	}
}

func (model *User) CollectionName() string {
	return "users"
}
