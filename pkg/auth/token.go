package auth

import (
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
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return a.hmacSecret, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("Invalid token")
	}

	expClaim, ok := claims["exp"]
	if !ok {
		return "", fmt.Errorf("Invalid expiration for token")
	}
	expInt, ok := expClaim.(int64)
	if !ok {
		return "", fmt.Errorf("Invalid expiration for token")
	}
	if time.Now().After(time.Unix(expInt, 0)) {
		return "", fmt.Errorf("Token expired")
	}

	idClaim, ok := claims["id"]
	if !ok {
		return "", fmt.Errorf("Invalid id for token")
	}
	id, ok := idClaim.(string)
	if !ok {
		return "", fmt.Errorf("Invalid id for token")
	}

	return id, nil
}
