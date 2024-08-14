package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	ID				primitive.ObjectID	`bson: "_id"`
	FirstName		*string				`json: "firstname" validate: "required,min=2,max=100"`
	LastName		*string				`json: "lastname" validate:"required min=2,max=100"`
	Password		*string				`json: "password" validate:"required min=6"`
	Email			*string				`json:"email" validate:"email,required"`
	Phone			*string				`json:"phone" validate:"required"`
	Token			*string				`json:"token"`
	UserType		*string				`json:"usertype" validate:"required,oneof=ADMIN USER"`
	RefreshToken	*string				`json:"refreshtoken"`
	CreatedAt		time.Time			`json:"createdat"`
	UpdatedAt		time.Time			`json:"updatedat"`
	UserId			string				`json:"userid"`
}


type Task struct {
	ID          string    `json:"id" bson:"id"` 
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}


type TaskRepository interface {
GetTasks(c context.Context) ([]*Task, error) 
	GetTasksById(c context.Context,id string) (*Task, error)
	 DeleteById(c context.Context,id string) (int64, error) 
	 UpdateTask(c context.Context,id string, updatedTask Task) error
	 AddTask(c context.Context,newTask Task) error
}

type TaskUsecase interface {
	GetTasks(c context.Context) ([]*Task, error) 
	GetTasksById(c context.Context,id string) (*Task, error)
	 DeleteById(c context.Context,id string) (int64, error) 
	 UpdateTask(c context.Context,id string, updatedTask Task) error
	 AddTask(c context.Context,newTask Task) error
	
}


type UserRepository interface {
    Signup(ctx context.Context, user User) (interface{}, error)
    Login(ctx context.Context, email string) (*User, error)
    GetUsers(ctx context.Context, startIndex int64, recordsPerPage int64) ([]*User, error)
    GetUser(ctx context.Context, user_id string) (User, error)
    Promote(ctx context.Context, user_id string, userType string) error
	UpdateAllTokens(token string, refreshToken string, user_id string) error

}
type UserUseCase interface {
    Signup(c context.Context, user User) (interface{}, error)
    Login(c context.Context, email string) (*User, error)
    GetUsers(c context.Context, startIndex int64, recordsPerPage int64) ([]*User, error)
    GetUser(c context.Context, user_id string) (User, error)
    Promote(c context.Context, user_id string, userType string) error
	UpdateAllTokens(token string, refreshToken string, user_id string) error

}