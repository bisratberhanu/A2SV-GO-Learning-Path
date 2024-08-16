package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"task_manger_clean_architecture/delivery/middleware"
	"task_manger_clean_architecture/domain"
	"task_manger_clean_architecture/domain/mocks"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TaskControllerTestSuite struct {
	suite.Suite
	router          *gin.Engine
	taskController  *TaskController
	mockTaskUseCase *mocks.TaskUsecase
	TaskGroup       []*domain.Task
	SingleTask      domain.Task
}


func (suite *TaskControllerTestSuite) SetupTest() {
  suite.router = gin.Default()
  suite.mockTaskUseCase = new(mocks.TaskUsecase)
  
  suite.taskController = &TaskController {
    TaskUseCase : suite.mockTaskUseCase,
  }

  suite.TaskGroup = []*domain.Task{
    {
      Title : "Title 1",
      Description : "this is title 1",
      DueDate : time.Now(),
      Status : "pending",
    },
    {
      Title : "Title 2",
      Description : "this is title 2",
      DueDate : time.Now(),
      Status : "pending",
    },
  }

  suite.SingleTask = domain.Task{Title : "Title 1", Description : "this is title 1",Status : "pending",}
  suite.router.GET("/task", middleware.Authenticate(), suite.taskController.GetTasks())
  suite.router.POST("/task", middleware.Authenticate(), suite.taskController.AddTask())
  suite.router.DELETE("/task/:task_id", middleware.Authenticate(), suite.taskController.DeleteById())
  suite.router.PUT("/task/:task_id", middleware.Authenticate(), suite.taskController.UpdateTask())
  suite.router.GET("/task/:task_id", middleware.Authenticate(), suite.taskController.GetTasksById())
}
type SignedDetails struct {
    Email     string
    FirstName string
    LastName  string
    UserType  string
    Uid       string
    jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")
func (suite *TaskControllerTestSuite) GenerateToken(email string, firstName string, lastName string, userType string, uid string) (string, error) {
    claims := &SignedDetails{
        Email:     email,
        FirstName: firstName,
        LastName:  lastName,
        UserType:  userType,
        Uid:       uid,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(),
        },
    }

    token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
    if err != nil {
        return "", errors.New("error while generating token")
    }

    return token, nil
}


func (suite *TaskControllerTestSuite) TestGetTasksSuccess() {
    // Create a new HTTP GET request to the /tasks endpoint without a request body
    returnedTasks := suite.TaskGroup
    req, err := http.NewRequest(http.MethodGet, "/task", nil)
    suite.NoError(err)

    // Generate a token for the user
    token, err := suite.GenerateToken("bisratbnegus@gmail.com","bisrat", "berhanu",  "user", "demoid", )
    suite.NoError(err)

    // Mock the GetTasks method to return the expected tasks
    suite.mockTaskUseCase.On("GetTasks", mock.Anything).Return(returnedTasks, nil)

    // Set the Content-Type and Authorization headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)

    // Create a new HTTP test recorder
    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    // Assert that the response status code is http.StatusOK
    suite.Equal(http.StatusOK, recorder.Code)

    // Unmarshal the response body to a slice of domain.Task
    var responseBody []*domain.Task
    fmt.Println(recorder.Body)
    err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
    suite.NoError(err)

    // Assert that the response body contains the expected tasks
    suite.Equal(returnedTasks[0].Title, responseBody[0].Title)
    suite.Equal(returnedTasks[0].Status, responseBody[0].Status)

    // Assert that the GetTasks method was called with the correct context
    suite.mockTaskUseCase.AssertCalled(suite.T(), "GetTasks", mock.Anything)
}


func (suite *TaskControllerTestSuite) TestGetTasks_ErrorFetchingTasks() {
    // Create a new HTTP GET request to the /tasks endpoint without a request body
    req, err := http.NewRequest(http.MethodGet, "/task", nil)
    suite.NoError(err)

    // Generate a token for the user
    token, err := suite.GenerateToken("bisratbnegus@gmail.com","bisrat", "berhanu",  "user", "demoid", )
    suite.NoError(err)

    // Mock the GetTasks method to return an error when called with any context
    suite.mockTaskUseCase.On("GetTasks", mock.Anything).Return(nil, errors.New("error while fetching tasks"))

    // Set the Content-Type and Authorization headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)

    // Create a new HTTP test recorder
    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    // Assert that the response status code is http.StatusInternalServerError
    suite.Equal(http.StatusInternalServerError, recorder.Code)

    // Unmarshal the response body to a gin.H map
    var responseBody gin.H
    err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
    suite.NoError(err)

    // Assert that the response body contains the correct error message
    suite.Equal("Failed to retrieve tasks", responseBody["error"])

    // Assert that the GetTasks method was called with any context
    suite.mockTaskUseCase.AssertCalled(suite.T(), "GetTasks", mock.Anything)
}




func (suite *TaskControllerTestSuite) TestAddTaskSuccess() {
    task := suite.SingleTask
    token, err := suite.GenerateToken("bisratbnegus@gmail.com","bisrat", "berhanu",  "ADMIN", "demoid", )
    suite.NoError(err)

    // Update mock to expect a context and the task as arguments
    suite.mockTaskUseCase.On("AddTask", mock.Anything, task).Return(nil)

    body, err := json.Marshal(task)
    suite.NoError(err, "no error while marshalling task data")

    req, err := http.NewRequest(http.MethodPost, "/task", bytes.NewBuffer(body))
    suite.NoError(err, "no error while creating new request")

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    // Expecting status OK (200)
    suite.Equal(http.StatusOK, recorder.Code)

    // Log the response body
    responseBodyBytes := recorder.Body.Bytes()
    fmt.Println("Response Body:", string(responseBodyBytes))

    var responseBody map[string]interface{}
    err = json.Unmarshal(responseBodyBytes, &responseBody)
    suite.NoError(err, "no error while unmarshalling response body")

    // Expecting the success message
    suite.Equal("Task added successfully", responseBody["message"])

    // Ensure the AddTask method was called with the correct arguments
    suite.mockTaskUseCase.AssertCalled(suite.T(), "AddTask", mock.Anything, task)
}


func (suite *TaskControllerTestSuite) TestDeleteByIdSuccess() {
    // Create a new DELETE request
    req, err := http.NewRequest(http.MethodDelete, "/task/12345", nil)
    suite.NoError(err)

    // Generate an authorization token
    token, err := suite.GenerateToken("bisratbnegus@gmail.com", "bisrat", "berhanu", "ADMIN", "demoid")
    suite.NoError(err)

    // Set up expectations for the mock task use case
    suite.mockTaskUseCase.On("DeleteById", mock.Anything, "12345").Return(int64(1), nil)

    // Set the request headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    // Record the response
    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    // Verify the response status code
    suite.Equal(http.StatusOK, recorder.Code)

    // Parse the response body
    var responseBody gin.H
    err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
    suite.NoError(err)

    // Verify the response message
    suite.Equal("Task deleted successfully", responseBody["message"])

    // Assert that the DeleteById method was called with the expected arguments
    suite.mockTaskUseCase.AssertCalled(suite.T(), "DeleteById", mock.Anything, "12345")
}


func (suite *TaskControllerTestSuite) TestUpdateTaskSuccess() {
    task := suite.SingleTask // Assuming SingleTask is a pre-defined task in your test suite
    
    // Generate an authorization token
    token, err := suite.GenerateToken("bisratbnegus@gmail.com", "bisrat", "berhanu", "ADMIN", "demoid")
    suite.NoError(err)

    // Set up expectations for the mock task use case
    suite.mockTaskUseCase.On("UpdateTask", mock.Anything, "12345", task).Return(nil)

    // Marshal the task data into JSON
    body, err := json.Marshal(task)
    suite.NoError(err, "no error while marshalling task data")

    // Create a new PUT request
    req, err := http.NewRequest(http.MethodPut, "/task/12345", bytes.NewBuffer(body))
    suite.NoError(err, "no error while creating new request")

    // Set the request headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    // Record the response
    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    // Verify the response status code
    suite.Equal(http.StatusOK, recorder.Code)

    // Log the response body
    responseBodyBytes := recorder.Body.Bytes()
    fmt.Println("Response Body:", string(responseBodyBytes))

    // Parse the response body
    var responseBody map[string]interface{}
    err = json.Unmarshal(responseBodyBytes, &responseBody)
    suite.NoError(err, "no error while unmarshalling response body")

    // Verify the response message
    suite.Equal(map[string]interface{}{"message": "Task updated successfully"}, responseBody)

    // Assert that the UpdateTask method was called with the expected arguments
    suite.mockTaskUseCase.AssertCalled(suite.T(), "UpdateTask", mock.Anything, "12345", task)
}

func (suite *TaskControllerTestSuite) TestGetTaskByIDSuccess() {
    // Create a new HTTP GET request to the /task/12345 endpoint without a request body
    returnedTask := suite.SingleTask // Assuming SingleTask is a pre-defined task in your test suite
    req, err := http.NewRequest(http.MethodGet, "/task/12345", nil)
    suite.NoError(err)

    // Generate a token for the user
    token, err := suite.GenerateToken("bisratbnegus@gmail.com", "bisrat", "berhanu", "ADMIN", "demoid")
    suite.NoError(err)

    // Mock the GetTasksById method to return the expected task
    suite.mockTaskUseCase.On("GetTasksById", mock.Anything, "12345").Return(&returnedTask, nil)

    // Set the Content-Type and Authorization headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    // Create a new HTTP test recorder
    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    // Assert that the response status code is http.StatusOK
    suite.Equal(http.StatusOK, recorder.Code)

    // Unmarshal the response body to a domain.Task object
    var responseBody domain.Task
    err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
    suite.NoError(err)

    // Assert that the response body contains the expected task data
    suite.Equal(returnedTask.Title, responseBody.Title)
    suite.Equal(returnedTask.Status, responseBody.Status)

    // Assert that the GetTasksById method was called with the correct parameters
    suite.mockTaskUseCase.AssertCalled(suite.T(), "GetTasksById", mock.Anything, "12345")
}


func TestControllerTestSuite(t *testing.T) {
  suite.Run(t, new(TaskControllerTestSuite))
}


