package routes

import (
	"github.com/bisratberhanu/A2SV-GO-Learning-Path/go_jwt_project/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("/signup", controllers.SignUp())
	incomingRoutes.POST("/login", controllers.Login())
}