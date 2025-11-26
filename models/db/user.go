package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

const (
	RoleUser     = "user"
	RoleAdmin    = "admin"
	RoleDoctor   = "doctor"
	RolePharmacy = "pharmacy"
	RoleLab      = "lab"
)

type User struct {
	gorm.Model
	Email         string    `json:"email" gorm:"uniqueIndex;not null"`
	Password      string    `json:"-" gorm:"not null"`
	Name          string    `json:"name" gorm:"not null"`
	Role          string    `json:"role" gorm:"not null;default:'user'"`
	EmailVerified bool      `json:"mail_verified" gorm:"column:email_verified;default:false"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Type  string `json:"type"`
}

func NewUser(email string, password string, name string, role string) *User {
	return &User{
		Email:         email,
		Password:      password,
		Name:          name,
		Role:          role,
		EmailVerified: false,
	}
}
