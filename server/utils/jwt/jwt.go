package jwt

import (
	"github.com/chenhqchn/ruohua/server/utils/config"
	"github.com/golang-jwt/jwt"
	"time"
)

type JwtClaims struct {
	Username string   `json:"username"`
	Role     []string `json:"role"`
	jwt.StandardClaims
}

// 生成 token
func GenToken(username string, role []string) (string, error) {
	claim := JwtClaims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(config.G().Jwt.ExpireTime)).Unix(),
			Issuer:    config.G().Jwt.Issuer,
		},
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := tokenObj.SignedString([]byte(config.G().Jwt.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析 token，判断是否合法、是否过期
func ParseToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(config.G().Jwt.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, err
	} else {
		if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, err
	}
}
