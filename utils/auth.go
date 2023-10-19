package utils

import (
	"Backend_TA/config"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var accesKey = config.RenderEnv("KEY_ACCES_TOKENS")
var secretAccesKey = []byte(accesKey)
var refreshKey = config.RenderEnv("KEY_REFRESH_TOKENS")
var secretRefreshKey = []byte(refreshKey)

func GenerateAccesTokens(claims *jwt.MapClaims) (string, error) {
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accesTokens, err := tokens.SignedString(secretAccesKey)
	if err != nil {
		return "", err
	}

	return accesTokens, nil
}

func GenerateRefreshTokens(claims *jwt.MapClaims) (string, error) {
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokens, err := tokens.SignedString(secretRefreshKey)
	if err != nil {
		return "", err
	}

	return refreshTokens, nil
}

func VerifyAccesToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return secretAccesKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func VerifyRefreshToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return secretRefreshKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func DecodeRefreshTokens(tokenString string) (jwt.MapClaims, error) {
	token, err := VerifyRefreshToken(tokenString)
	if err != nil {
		return nil, err
	}

	//isOk merupakan variabel yang menampung data apakah valid atau tidak
	claims, isOk := token.Claims.(jwt.MapClaims)
	if isOk && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
