package routers

import (
	"task_manger_clean_architecture/delivery/controllers"
	"task_manger_clean_architecture/repositories"
	"task_manger_clean_architecture/usecases"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewTaskRouter( timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	tr := repositories.NewTaskRepository(db, "task")
	tc := &controllers.TaskController{
		TaskUseCase: usecases.NewTaskUseCase(tr, timeout),
	}
	group.POST("/task", tc.AddTask())
	group.GET("/task", tc.GetTasks())
	group.GET("/task/:task_id", tc.GetTasksById())
	group.DELETE("/task/:task_id", tc.DeleteById())
	group.PUT("/task/:task_id", tc.UpdateTask())
}