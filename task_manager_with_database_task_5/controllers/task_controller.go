package controllers

import ("github.com/gin-gonic/gin"
"task_manager_with_database_task_5/data"
)

func GetTasks(c *gin.Context) {
	data.GetTasksDb(c)

}

func GetTasksById(c *gin.Context) {
	id:= c.Param("id")
	data.GetTasksByIdDb(id)
}

func UpdateTask(c *gin.Context) {
	


}

func DeleteById(c *gin.Context) {
	id:= c.Param("id")
	data.DeleteByIdDb(id)

}

func AddTask(c *gin.Context) {

}