package main

import (
	"task_manager_with_database_task_5/router"

	"github.com/gin-gonic/gin"
)

func main(){
	real_Router := gin.Default()
	router.Router(real_Router)

}