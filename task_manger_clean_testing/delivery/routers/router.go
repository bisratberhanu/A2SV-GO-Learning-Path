package routers

import (
	"task_manger_clean_architecture/delivery/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(timeout time.Duration, db *mongo.Database, gin *gin.Engine){
	publicRouter:= gin.Group("")
	NewSignUPRouter( timeout,db, publicRouter)
	NewLoginRouter(timeout,db,publicRouter)
	protectedRouter:= gin.Group("")
	protectedRouter.Use(middleware.Authenticate())
	NewUserRouter(timeout, db, protectedRouter)
	
	NewTaskRouter(timeout, db, protectedRouter)
} 