package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//GetNewToken generates a new token for the given id
func (a *Authenticator) GetNewToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // the token is valid for 30 days
	})

	return token.SignedString(a.hmacSecret)
}

//ValidateToken validates the token and returns the user ID if it's valid
func (a *Authenticator) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Validating signing method (alg)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return a.hmacSecret, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	idClaim, ok := claims["id"]
	if !ok {
		return "", errors.New("invalid id for token")
	}

	id, ok := idClaim.(string)

	if !ok {
		return "", errors.New("invalid id for token")
	}

	return id, nil
}
