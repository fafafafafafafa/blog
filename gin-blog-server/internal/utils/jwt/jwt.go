package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type MyClaims struct {
	UserId  int   `json:"user_id"`
	RoleIds []int `json:"role_ids"`
	// UUID    string `json:"uuid"`
	jwt.RegisteredClaims
}

func GenToken(secret, issuer string, expireHour, userId int, roleIds []int) (string, error) {
	claims := MyClaims{
		UserId:  userId,
		RoleIds: roleIds,
		// UUID:    uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
