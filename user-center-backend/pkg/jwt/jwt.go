package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"user-center-backend/constant"
	"user-center-backend/pkg/errcode"
)

type BlogClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return constant.JwtSecret, nil
}

const AccessTokenExpireDuration = time.Hour * 24
const RefreshTokenExpireDuration = time.Hour * 24 * 7

func GenToken(userID uint64, username string) (accessToken, refreshToken string, err error) {
	c := BlogClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenExpireDuration).Unix(),
			Issuer:    "user-center-backend",
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(constant.JwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RefreshTokenExpireDuration).Unix(),
		Issuer:    "user-center-backend",
	}).SignedString(constant.JwtSecret)
	if err != nil {
		return "", "", err
	}

	return
}

func ParseToken(tokenString string) (claims *BlogClaims, err error) {
	var token *jwt.Token
	claims = new(BlogClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return
	}
	if !token.Valid {
		err = errcode.InvalidToken
	}
	return
}
