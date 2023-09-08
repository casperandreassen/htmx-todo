package utils

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/ed25519"
	"time"
)

var privateKey ed25519.PrivateKey
var publicKey ed25519.PublicKey

func GenerateKeyPairs() bool {
	pub, pri, err := ed25519.GenerateKey(rand.Reader)

	if err == nil {
		privateKey = pri
		publicKey = pub
		return true
	}
	return false
}

func IssueJwtToken(id int, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"iss": "ruth",
		"exp": time.Now().Add(time.Hour * 8765).UnixMilli(),
		"data": map[string]interface{}{
			"id":       id,
			"username": username,
		},
	})
	tokenString, err := token.SignedString(privateKey)
	return tokenString, err
}

func VerifyJwtToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		ok := token.Method.Alg() == "EdDSA"
		if ok {
			return publicKey, nil
		} else {
			return "", nil
		}
	})
	if err != nil {
		return nil, errors.New("token not valid")
	} else {
		return parsedToken, nil
	}
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
