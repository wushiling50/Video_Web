package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWTsecret = []byte("ABAB")

type Claims struct {
	Id       uint   `json:"id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// 签发token
func GenerateToken(id uint, username, password string) (string, error) {
	notTime := time.Now()
	expireTime := notTime.Add(24 * time.Hour)
	claims := Claims{
		Id:       id,
		UserName: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "video_web",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(JWTsecret)
	return token, err
}

// 验证token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTsecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

type EmailClaims struct {
	Uid       uint   `json:"uid"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Operation uint   `json:"operation"`
	jwt.StandardClaims
}

// 签发emailtoken
func GenerateEmailToken(uid, operation uint, email, password string) (string, error) {
	notTime := time.Now()
	expireTime := notTime.Add(24 * time.Hour)
	claims := EmailClaims{
		Uid:       uid,
		Email:     email,
		Password:  password,
		Operation: operation,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "video_web",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(JWTsecret)
	return token, err
}

// 验证emailtoken
func ParseEmailToken(token string) (*EmailClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTsecret, nil
	})
	if tokenClaims != nil {
		if emailClaims, ok := tokenClaims.Claims.(*EmailClaims); ok && tokenClaims.Valid {
			return emailClaims, nil
		}
	}
	return nil, err
}
