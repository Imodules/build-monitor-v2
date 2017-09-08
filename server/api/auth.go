package api

import (
	"errors"

	"build-monitor-v2/server/db"

	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTClaims struct {
	jwt.StandardClaims
	UserId   string   `json:"userId"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Groups   []string `json:"groups"`
}

var TokenIsMissingRequiredField = errors.New("Token is missing a required field")

func GenerateToken(user *db.User, secret string) (string, error) {

	// we have a good payloadFactory ~ generate a token
	claims := &JWTClaims{
		UserId:   user.Id.Hex(),
		Username: user.Username,
		Email:    user.Email,
	}
	claims.ExpiresAt = time.Now().Add(time.Minute * 60).Unix()

	// create a token with our secret key
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (c JWTClaims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	if c.UserId == "" {
		return TokenIsMissingRequiredField
	}

	return nil
}
