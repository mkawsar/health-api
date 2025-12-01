package services

import (
	"errors"
	db "health/models/db"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateToken creates a new JWT token for the given user, with the given type and expiration time.
// The token is signed with the ES256 algorithm and the secret key from the .env file.
// The token is then saved to the tokens collection in the database.
// If the token cannot be created or saved, an error is returned.
func CreateToken(user *db.User, tokenType string, expiresAt time.Time) (*db.Token, error) {
	claims := &db.UserClaims{
		Email: user.Email,
		Type:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(Config.JWTSecretKey))
	if err != nil {
		return nil, errors.New("cannot create access token")
	}

	tokenModel := db.NewToken(user.ID, tokenString, tokenType, expiresAt)
	err = mgm.Coll(tokenModel).Create(tokenModel)
	if err != nil {
		return nil, errors.New("cannot save access token to db")
	}

	return tokenModel, nil
}

// DeleteTokenById deletes a token from the database by its ID.
// If the token does not exist or the deletion fails, an error is returned.
func DeleteTokenById(id primitive.ObjectID) error {
	ctx := mgm.Ctx()
	deleteResult, err := mgm.Coll(&db.Token{}).DeleteOne(ctx, bson.M{field.ID: id})

	if err != nil || deleteResult.DeletedCount <= 0 {
		return errors.New("cannot delete token")
	}

	return nil
}

// GenerateAccessTokens creates a new access and refresh token for the given user.
// The access token has a TTL of JWT_ACCESS_EXPIRATION_MINUTES minutes and the refresh token has a TTL of JWT_REFRESH_EXPIRATION_DAYS days.
// The tokens are signed with the ES256 algorithm and the secret key from the .env file.
// The tokens are then saved to the tokens collection in the database.
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

	tokenModel := &db.Token{}
	userId, _ := primitive.ObjectIDFromHex(claims.Subject)
	err = mgm.Coll(tokenModel).First(
		bson.M{"type": tokenType, "user": userId, "blacklisted": false},
		tokenModel,
	)
	if err != nil {
		return nil, errors.New("cannot find token")
	}

	return tokenModel, nil
}
