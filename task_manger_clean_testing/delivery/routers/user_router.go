package routers

import (
	"task_manger_clean_architecture/delivery/controllers"
	"task_manger_clean_architecture/repositories"
	"task_manger_clean_architecture/usecases"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRouter(timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ur := repositories.NewUserRepository(db, "user")
	uc := &controllers.UserController{
		UserUseCase: usecases.NewUserUseCase(ur, timeout),
	}
	group.GET("/users", uc.GetUsers())
	group.GET("/users/:user_id", uc.GetUser())
	group.POST("/promote/:user_id", uc.Promote())
}