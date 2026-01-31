package api

import "github.com/gin-gonic/gin"

func errError(err error) gin.H {
	return gin.H{
		"err": gin.H{
			"body": err.Error(),
		},
	}
}

func errString(str string) gin.H {
	return gin.H{
		"err": gin.H{
			"body": str,
		},
	}
}
