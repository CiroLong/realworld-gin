package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtManager struct {
	secret     []byte
	expireTime time.Duration
}

func NewManager(secret string, expire time.Duration) Manager {
	return &jwtManager{
		secret:     []byte(secret),
		expireTime: expire,
	}
}

func (m *jwtManager) Generate(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *jwtManager) Parse(tokenStr string) (int64, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return m.secret, nil
		},
	)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*Claims) // 强转为预定义的Claims
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}
