package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Token struct {
	gorm.Model
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	Token       string    `json:"token" gorm:"not null;index"`
	Type        string    `json:"type" gorm:"not null"`
	ExpriesAt   time.Time `json:"expries_at" gorm:"not null;index"`
	BlackListed bool      `json:"blacklisted" gorm:"column:blacklisted;default:false"`
}

// GetResponseJson returns a gin.H representation of the token that is safe for transmission over the network.
// It contains the token string and the expiration time in the "expires" field, formatted as "2006-01-02 15:04:05".
func (model *Token) GetResponseJson() gin.H {
	return gin.H{"token": model.Token, "expires": model.ExpriesAt.Format("2025-02-04 15:04:05")}
}

// NewToken creates a new Token with the given user id, token string, token type and expiration time.
func NewToken(userId uint, tokenString string, tokenType string, expriesAt time.Time) *Token {
	return &Token{
		UserID:      userId,
		Token:       tokenString,
		Type:        tokenType,
		ExpriesAt:   expriesAt,
		BlackListed: false,
	}
}
