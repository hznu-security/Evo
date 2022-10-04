/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/09 16:32
 * 描述     ：jwt生成token和解析部分
 */

package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

var jwtKey = []byte("a_secret_crect")

// TEAM 两种角色，0代表是队伍，1代表是管理员
const TEAM uint8 = 0
const ADMIN uint8 = 1

type Claims struct {
	ID   uint
	Role uint8
	jwt.StandardClaims
}

func ReleaseToken(id uint, role uint8) (string, error) {
	validity := viper.GetDuration("auth.vilidity")
	expirationTime := time.Now().Add(validity * time.Hour)
	if role != TEAM && role != ADMIN {
		return "", errors.New("角色类型不正确")
	}
	claims := &Claims{
		ID:   id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "evo",
			Subject:   "token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (i interface{}, err error) {
			return jwtKey, nil
		})
	return token, claims, err
}
