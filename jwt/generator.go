package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Generator struct {
	SigningKey string
}

func (g Generator) Generate(ownerUUID string) (string, error) {
	claims := StoryscriptClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "storyscript",
			IssuedAt:  time.Now().UTC().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 365).UTC().Unix(),
		},
		OwnerUUID: ownerUUID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(g.SigningKey))

	return tokenString, nil
}

type StoryscriptClaims struct {
	jwt.StandardClaims
	OwnerUUID string `json:"owner_uuid"`
}
