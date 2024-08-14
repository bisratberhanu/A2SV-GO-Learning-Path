package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"task_manger_clean_architecture/domain"
	"task_manger_clean_architecture/infrastructure"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
var validate = validator.New()
type UserController struct {
	UserUseCase domain.UserUseCase
}

func (uc *UserController) Signup() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        var user domain.User
        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        validationError := validate.Struct(user)
        if validationError != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
            return
        }

        // Hash the password
        password := infrastructure.HashPassword(*user.Password)
        user.Password = &password
        // Set additional user fields
        user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        user.ID = primitive.NewObjectID()
        user.UserId = user.ID.Hex()

        // Call the use case to sign up the user
        resultInsertionNumber, err := uc.UserUseCase.Signup(ctx, user)
        if err != nil {
            msg := fmt.Sprint("user item was not created")
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
        }

        c.JSON(http.StatusOK, gin.H{"insertionnumber": resultInsertionNumber})
    }
}

func (uc *UserController) Login() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        var user domain.User
        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Call the use case to login the user
        foundUser, err := uc.UserUseCase.Login(ctx, *user.Email)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password not found"})
            return
        }

        // Verify the password
        passwordIsValid, msg := infrastructure.VerifyPassword(*user.Password, *foundUser.Password)
        if !passwordIsValid {
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
        }

        if foundUser.Email == nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "user not found"})
            return
        }

        // Generate tokens
        token, refreshToken, _ := infrastructure.GenerateAllTokens(*foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.UserType, &foundUser.UserId)
        uc.UserUseCase.UpdateAllTokens(token, refreshToken, foundUser.UserId)

        c.JSON(http.StatusOK, foundUser)
    }
}

func (uc *UserController) GetUsers() gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := infrastructure.CheckUserType(c, "ADMIN"); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        recordsPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
        if err != nil || recordsPerPage < 1 {
            recordsPerPage = 10
        }

        page, err := strconv.Atoi(c.Query("page"))
        if err != nil || page < 1 {
            page = 1
        }

        startIndex := int64((page - 1) * recordsPerPage)

        users, err := uc.UserUseCase.GetUsers(ctx, startIndex, int64(recordsPerPage))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing user items"})
            return
        }

        c.JSON(http.StatusOK, users)
    }
}

func (uc *UserController) GetUser() gin.HandlerFunc {
    return func(c *gin.Context) {
        userId := c.Param("user_id")
        if err := infrastructure.MatchUserTypeToUid(c, userId); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        user, err := uc.UserUseCase.GetUser(ctx, userId)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, user)
    }
}

func (uc *UserController) Promote() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check if the user is an admin
        userType, exists := c.Get("usertype")
        if !exists || userType != "ADMIN" {
            c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
            return
        }

        // Get the user ID from the request
        userId := c.Param("user_id")

        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        // Call the use case to promote the user
        err := uc.UserUseCase.Promote(ctx, userId, "ADMIN")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user type"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin successfully"})
    }
}

