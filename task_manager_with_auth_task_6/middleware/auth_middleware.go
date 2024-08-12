package middleware

import (
	"net/http"
	"task_manager_with_auth_task_6/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc{
	return func(c *gin.Context){
		clientToken:= c.Request.Header.Get("token")
		if clientToken==""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not Authorized"})
			c.Abort()
			return
			
		}
		claims,err := helpers.ValidateToken(clientToken)
		if err!=""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
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

