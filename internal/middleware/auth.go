package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"github/CiroLong/realworld-gin/internal/common"
	"github/CiroLong/realworld-gin/internal/model"
	"net/http"
	"strings"
)

// Strips 'TOKEN ' prefix from token string
func stripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 5 && strings.ToUpper(tok[0:6]) == "TOKEN " {
		return tok[6:], nil
	}
	return tok, nil
}

// AuthorizationHeaderExtractor
//
//	Extract  token from Authorization header
//	Uses PostExtractionFilter to strip "TOKEN " prefix from header
var AuthorizationHeaderExtractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{"Authorization"},
	Filter:    stripBearerPrefixFromTokenString,
}

// MyAuth2Extractor
//
//	Extractor for OAuth2 access tokens.  Looks in 'Authorization'
//	header then 'access_token' argument for a token.
var MyAuth2Extractor = &request.MultiExtractor{
	AuthorizationHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

// UpdateContextUserModel A helper to write user_id and user_model to the context
func UpdateContextUserModel(c *gin.Context, myUserId uint) {
	var myUserModel model.UserModel
	if myUserId != 0 {
		db := common.GetDB()
		db.First(&myUserModel, myUserId)
	}
	c.Set("my_user_id", myUserId)
	c.Set("my_user_model", myUserModel)
}

// AuthMiddleware Auth middleware parses JWT, sets current user in context
func AuthMiddleware(auto401 bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		UpdateContextUserModel(c, 0)
		token, err := request.ParseFromRequest(c.Request, MyAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			b := []byte(common.NBSecretPassword)
			return b, nil
		})
		if err != nil {
			if auto401 {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			myUserId := uint(claims["id"].(float64))
			//fmt.Println(my_user_id,claims["id"])
			UpdateContextUserModel(c, myUserId)
		}
	}
}
