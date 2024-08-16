package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"task_manger_clean_architecture/infrastructure"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not Authorized"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix from the token string
		tokenParts := strings.Split(clientToken, " ")
		if len(tokenParts) == 2 {
			clientToken = tokenParts[1]
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}

		claims, err := infrastructure.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"errorsss": err})
			fmt.Println("check error happend")
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("firstname", claims.FirstName)
		c.Set("lastname", claims.LastName)
		c.Set("uid", claims.Uid)
		c.Set("usertype", claims.UserType)
	}
}
