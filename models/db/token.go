package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Token struct {
	mgm.DefaultModel `bson:",inline"`
	User             primitive.ObjectID `json:"user" bson:"user"`
	Token            string             `json:"token" bson:"token"`
	Type             string             `json:"type" bson:"type"`
	ExpriesAt        time.Time          `json:"expries_at" bson:"expries_at"`
	BlackListed      bool               `json:"blacklisted" bson:"blacklisted"`
}

// GetResponseJson returns a gin.H representation of the token that is safe for transmission over the network.
// It contains the token string and the expiration time in the "expires" field, formatted as "2006-01-02 15:04:05".
func (model *Token) GetResponseJson() gin.H {
	return gin.H{"token": model.Token, "expires": model.ExpriesAt.Format("2025-02-04 15:04:05")}
}

// NewToken creates a new Token with the given user id, token string, token type and expiration time.
func NewToken(userId primitive.ObjectID, tokenString string, tokenType string, expriesAt time.Time) *Token {
	return &Token{
		User:        userId,
		Token:       tokenString,
		Type:        tokenType,
		ExpriesAt:   expriesAt,
		BlackListed: false,
	}
}

// CollectionName returns the name of the collection that stores Token documents.
func (model *Token) CollectionName() string {
	return "tokens"
}
