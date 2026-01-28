package jwt

import "github.com/golang-jwt/jwt/v5"

// 这里的Claims 的 claim是jwt中的概念
type Claims struct {
	UserID int64 `json:"uid"`
	jwt.RegisteredClaims
}
