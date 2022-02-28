package utils

import (
	"github.com/dgrijalva/jwt-go"
)

var secretKey = "sdafsadgdfsagbfsdbdsfbsdfasdf23423534hui[]=230-3hdsj"

func GenerateToken(id uint, email string) string {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := parseToken.SignedString([]byte(secretKey))

	return signedToken
}
