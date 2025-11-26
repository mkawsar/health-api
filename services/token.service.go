package services

import (
	"errors"
	"fmt"
	db "health/models/db"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// CreateToken creates a new JWT token for the given user, with the given type and expiration time.
// The token is signed with the HS256 algorithm and the secret key from the .env file.
// The token is then saved to the tokens table in the database.
// If the token cannot be created or saved, an error is returned.
func CreateToken(user *db.User, tokenType string, expiresAt time.Time) (*db.Token, error) {
	claims := &db.UserClaims{
		Email: user.Email,
		Type:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(Config.JWTSecretKey))
	if err != nil {
		return nil, errors.New("cannot create access token")
	}

	tokenModel := db.NewToken(user.ID, tokenString, tokenType, expiresAt)
	err = DB.Create(tokenModel).Error
	if err != nil {
		return nil, errors.New("cannot save access token to db")
	}

	return tokenModel, nil
}

// DeleteTokenById deletes a token from the database by its ID.
// If the token does not exist or the deletion fails, an error is returned.
func DeleteTokenById(id uint) error {
	result := DB.Delete(&db.Token{}, id)
	if result.Error != nil || result.RowsAffected == 0 {
		return errors.New("cannot delete token")
	}

	return nil
}

// GenerateAccessTokens creates a new access and refresh token for the given user.
// The access token has a TTL of JWT_ACCESS_EXPIRATION_MINUTES minutes and the refresh token has a TTL of JWT_REFRESH_EXPIRATION_DAYS days.
// The tokens are signed with the HS256 algorithm and the secret key from the .env file.
// The tokens are then saved to the tokens table in the database.
// If either token cannot be created or saved, an error is returned.
func GenerateAccessTokens(user *db.User) (*db.Token, *db.Token, error) {
	accessExpiresAt := time.Now().Add(time.Duration(Config.JWTAccessExpirationMinutes) * time.Minute)
	refreshExpiresAt := time.Now().Add(time.Duration(Config.JWTRefreshExpirationDays) * time.Hour * 24)

	accessToken, err := CreateToken(user, db.TokenTypeAccess, accessExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := CreateToken(user, db.TokenTypeRefresh, refreshExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func VerifyToken(token string, tokenType string) (*db.Token, error) {
	claims := &db.UserClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.JWTSecretKey), nil
	})

	if err != nil || claims.Type != tokenType {
		return nil, errors.New("not valid token")
	}

	if time.Since(claims.ExpiresAt.Time) > 10*time.Second {
		return nil, errors.New("token is expired")
	}

	userId, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user id in token")
	}

	tokenModel := &db.Token{}
	err = DB.Where("type = ? AND user_id = ? AND blacklisted = ?", tokenType, uint(userId), false).First(tokenModel).Error
	if err != nil {
		return nil, errors.New("cannot find token")
	}

	return tokenModel, nil
}
