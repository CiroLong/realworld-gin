package middleware

// Auth Middleware 负责从 HTTP 请求中解析 JWT，校验合法性，并将用户身份注入上下文。
// internal/middleware/auth.go
import (
	"github/CiroLong/realworld-gin/internal-v2/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtMgr jwt.Manager) gin.HandlerFunc {
	// 注意Middleware写法：一个闭包函数
	return func(c *gin.Context) {

		// 1. 取 Authorization Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": gin.H{
					"body": []string{"missing authorization header"},
				},
			})
			return
		}

		// 2. Bearer token 解析
		// HTTP头部中token字段约定
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": gin.H{
					"body": []string{"invalid authorization header"},
				},
			})
			return
		}

		token := parts[1]

		// 3. 校验并解析 JWT
		userID, err := jwtMgr.Parse(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": gin.H{
					"body": []string{"invalid or expired token"},
				},
			})
			return
		}

		// 4. 注入上下文
		// 这里将userID注入context中向下传递
		c.Set(ContextUserIDKey, userID)

		// 5. 放行
		c.Next()
	}
}

// OptionalAuthMiddleware 这个中间件是考虑到部分endpoint对带token和不带token的请求响应不同
func OptionalAuthMiddleware(jwtMgr jwt.Manager) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Token" {
			c.Next()
			return
		}

		userID, err := jwtMgr.Parse(parts[1])
		if err == nil {
			c.Set(ContextUserIDKey, userID)
		}

		c.Next()
	}
}
