package utils

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt"
	"go-server/db"
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

func IssueJwtToken(user db.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"iss": "ruth",
		"exp": time.Now().Add(time.Hour * 8765).UnixMilli(),
		"data": map[string]interface{}{
			"id":       user.Id,
			"username": user.Username,
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

func TransformTodos(todos []db.DBTodo) ([]db.DBTodo, []db.DBTodo, []db.DBTodo, error) {
	completedTodos := []db.DBTodo{}
	expiredTodos := []db.DBTodo{}
	otherTodos := []db.DBTodo{}
	for i := range todos {
		if !todos[i].Date.Valid {
			if todos[i].Status == 1 {
				completedTodos = append(completedTodos, todos[i])
			} else {
				otherTodos = append(otherTodos, todos[i])
			}
		} else {
			todoDate, err := time.Parse("2006-01-02", todos[i].Date.String)
			if err != nil {
				return nil, nil, nil, errors.New("Could not parse date")
			}
			if todos[i].Status == 1 {
				completedTodos = append(completedTodos, todos[i])
			} else {
				if time.Now().After(todoDate) {
					expiredTodos = append(expiredTodos, todos[i])
				} else {
					otherTodos = append(otherTodos, todos[i])
				}
			}
		}

	}
	return completedTodos, expiredTodos, otherTodos, nil
}
