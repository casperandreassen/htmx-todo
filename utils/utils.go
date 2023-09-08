package utils

import (
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type TodoClaims struct {
	Userid   int    `json:"userid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtSecret = os.Getenv("JWT_SECRET")
var jwtKey = []byte(jwtSecret)

func IssueJwtToken(id int, username string) (string, error) {
	claims := TodoClaims{
		id,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8765 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ruth",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func VerifyJwtToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TodoClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, errors.New("Could not parse token.")
	}

	return token, nil
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
