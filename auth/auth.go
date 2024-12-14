package auth

import (
	"bytes"
	"crypto/sha256"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var Domain string
var SigningKey string

func CheckPassword(serverPass, clientPass []byte) error {
	return bcrypt.CompareHashAndPassword(serverPass, clientPass)
}

func CreateJWT(username string) (string, error) {

	var aud []string
	aud = append(aud, Domain)

	id := make([][]byte, 2)
	id[0] = []byte(username)
	id[1] = []byte(Domain)
	idString := sha256.Sum224(bytes.Join(id, []byte("")))

	claims := jwt.RegisteredClaims{
		Issuer:    Domain,
		Subject:   username,
		Audience:  aud,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        string(idString[:]),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(SigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func CheckToken(tokenString string) error {

	claims := &jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return SigningKey, nil
	})
	if err != nil {
		return err
	}
	return nil
}
