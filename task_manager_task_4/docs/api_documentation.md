## Introduction
The Task Manager is a simple RESTful API built using the Go programming language and the Gin framework. It provides functionalities to manage tasks including creating, reading, updating, and deleting tasks. This documentation provides an overview of the project structure and explains the functionality of each file in the project.

## Project Structure
The project consists of the following main components:

**controllers/**

**data/**

**models/**

**routers/**

**main.go** 

1. ### controllers/task_controller.go

Purpose: This file initializes the task controller.

Explanation: The StartTaskController function initializes the routing by calling the Router function from the routers package.

2. ### data/task_services.go
**Purpose**: This file contains the service functions for handling task-related operations such as fetching, updating, deleting, and adding tasks.

**Explanation**:

**GetTasks**(c *gin.Context): Retrieves all tasks and responds with a JSON array of tasks.

**GetTasksById**(c *gin.Context): Retrieves a task by its ID. If the task is found, it responds with the task's details; otherwise, it responds with a "task not found" message.

**UpdateTask**(c *gin.Context): Updates an existing task by its ID. If the task is found and updated successfully, it responds with a success message; otherwise, it responds with an appropriate error message.

**DeleteById**(c *gin.Context): Deletes a task by its ID. If the task is found and deleted successfully, it responds with a success message; otherwise, it responds with a "ID not found" message.

**AddTask**(c *gin.Context): Adds a new task to the list of tasks. If the task is added successfully, it responds with a success message; otherwise, it responds with an error message.

3. ### models/task.go
Purpose: This file defines the data model for the task and provides mock data for initial testing.

Explanation:

**Task struct**: Defines the structure of a task with fields ID, Title, Description, DueDate, and Status.
**Tasks**: A slice of Task structs initialized with some mock data to simulate a database.

4. ### routers/router.go
**Purpose**: This file sets up the routing for the application using the Gin framework.

Explanation:

**Router()**: Initializes the router, sets up the endpoints, and associates each endpoint with its corresponding handler function from the data package. The endpoints include:

**GET /tasks**: Retrieves all tasks.

**GET /tasks/:id**: Retrieves a task by its ID.

**PUT /tasks/:id**: Updates a task by its ID.

**DELETE /tasks/:id**: Deletes a task by its ID.

**POST /tasks**: Adds a new task.

**The router runs on localhost:8080.**

5. ## main.go
**Purpose**: The entry point of the application.


The postman documentation can be found [here](https://documenter.getpostman.com/view/37520949/2sA3rzKt1E)