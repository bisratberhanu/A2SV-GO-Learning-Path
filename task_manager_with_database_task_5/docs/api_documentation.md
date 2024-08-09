## Folder Structure
```
task_manager/
├── main.go
├── controllers/
│   └── task_controller.go
├── models/
│   └── task.go
├── data/
│   └── task_service.go
├── router/
│   └── router.go
├── docs/
│   └── api_documentation.md
└── go.mod

```
1. ## main.go
**Purpose**: The entry point of the application. It initializes the server, sets up the router, and starts the application.

2. ## controllers/
**task_controller.go:**

**Purpose**: Contains the logic for handling HTTP requests related to tasks. It uses the functions defined in the data package to interact with the database. The controller is responsible for sending appropriate responses to the client.

3. ## models/
**task.go:**

**Purpose**: Defines the Task struct which represents the schema of the task documents stored in MongoDB. This model is used throughout the application to ensure consistent data handling.
4. ## data/

**task_service.go:**

**Purpose**: Implements the CRUD operations to interact with the MongoDB database. This package is responsible for directly communicating with the database, performing queries, and returning results.
API Functions:

**GetTasksDb**:

**Description**: Retrieves all tasks from the tasks collection in MongoDB.
Error Handling: Returns an error if the database query fails or if there's an issue decoding the task data.
Return: A slice of task pointers and an error (if any).

**GetTasksByIdDb**:

**Description**: Fetches a single task by its id from the tasks collection.
Error Handling: Returns an error if the task is not found or if there's an issue with the query.
Return: A task pointer and an error (if any).

**DeleteByIdDb**:

**Description**: Deletes tasks by their id.
Error Handling: Returns the number of deleted documents and an error if the deletion fails.
Return: The number of documents deleted and an error (if any).

**UpdateTaskDb**:

**Description**: Updates a task's fields based on the provided id.
Error Handling: Returns an error if the update fails or if no fields are provided for the update.
Return: An error (if any).

**AddTaskDb**:

Description: Adds a new task to the tasks collection.
Error Handling: Returns an error if the insertion fails.
Return: An error (if any).

## Database Initialization:
The database connection is initialized in the Db function within the data package. This function is responsible for connecting to the MongoDB instance using the provided URI. Once connected, the function pings the database to ensure the connection is active.

### Db():

**Description**: Connects to the MongoDB database using the MongoDB Go Driver.

**Steps**:
Client Options: Uses mongo/options.ClientOptions to set the URI for the MongoDB cluster.

**Connect**: Attempts to connect to MongoDB using mongo.Connect.

**Ping**: Verifies the connection by pinging the database.
Return: Returns a handle to the task_management database.
### DbTasks:

**Description**: A global variable that holds the reference to the tasks collection within the task_management database. This variable is used across the task_service.go file for CRUD operations.

5. ## router/
**router.go:**

**Purpose**: Defines the application's routes and associates them with the appropriate controller functions. It sets up the API endpoints for task-related operations.
## 6. docs/

**api_documentation.md:**
Purpose: Contains the API documentation for the project. This file explains how to interact with the API, including available endpoints, request parameters, and response formats.
7. ## go.mod

**Purpose**: The Go module file which manages the project's dependencies


the full API documentation can be found here: https://documenter.getpostman.com/view/37520949/2sA3s1os4C

