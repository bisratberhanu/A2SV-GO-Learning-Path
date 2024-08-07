package data
import(
	"github.com/gin-gonic/gin"
	"net/http"
	"task_manager_task_4/models"

)

func GetTasks(c *gin.Context){
	c.IndentedJSON(http.StatusOK, models.Tasks )
}


func GetTasksById(c *gin.Context){
	id:= c.Param("id")
	for _, value := range models.Tasks{
		if value.ID == id{
			c.IndentedJSON(http.StatusOK, value)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message":"task not found"})

}


func UpdateTask(c *gin.Context){
	id:= c.Param("id")
	var updatedTask models.Task

	if err:= c.ShouldBindJSON(&updatedTask); err!= nil{
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	for idx, task := range models.Tasks{
		if task.ID == id{
			if updatedTask.Title!= ""{
				models.Tasks[idx].Title = updatedTask.Title
			}
			if updatedTask.Description!= ""{
				models.Tasks[idx].Description = updatedTask.Description
			}

			c.IndentedJSON(http.StatusOK, gin.H{"message": "task updated "})
			return
			
			
	}
}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "id not found"})
}


func DeleteById(c *gin.Context){
	id:= c.Param("id")
	for idx, task:= range models.Tasks{
		if task.ID == id{
			models.Tasks = append(models.Tasks[:idx], models.Tasks[idx+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
			return
		}

	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "ID not found"})
}


func AddTask(c *gin.Context){
	var newTask models.Task
	if err:= c.BindJSON(&newTask); err!= nil{
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	models.Tasks = append(models.Tasks, newTask)
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "task created"})
}