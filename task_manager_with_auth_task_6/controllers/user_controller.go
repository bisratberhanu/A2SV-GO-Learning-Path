package controllers

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"

    "task_manager_with_auth_task_6/database"
    "task_manager_with_auth_task_6/helpers"
    "task_manager_with_auth_task_6/models"
    "github.com/go-playground/validator/v10"
    "golang.org/x/crypto/bcrypt"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var userCollections *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        log.Panic(err)
    }
    return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
    err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
    check := true
    msg := ""
    if err != nil {
        check = false
        msg = "password or email not found"
    }
    return check, msg
}

func Register() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        var user models.User
        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        validationError := validate.Struct(user)
        if validationError != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
            return
        }
        password := HashPassword(*user.Password)
        user.Password = &password
        count, err := userCollections.CountDocuments(ctx, bson.M{"email": user.Email})
        defer cancel()
        if err != nil {
            log.Panic(err)
            c.JSON(http.StatusBadRequest, gin.H{"error": "error occurred while getting email"})
        }
        count, err = userCollections.CountDocuments(ctx, bson.M{"phone": user.Phone})
        defer cancel()
        if err != nil {
            log.Panic(err)
            c.JSON(http.StatusBadRequest, gin.H{"error": "error occurred while getting phone"})
            return
        }
        if count > 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "email or phone number already exists"})
            return
        }
        user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        user.ID = primitive.NewObjectID()
        user.UserId = user.ID.Hex()
        token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, user.FirstName, user.LastName, user.UserType, &user.UserId)
        user.Token = &token
        user.RefreshToken = &refreshToken
        resultInsertionNumber, insertErr := userCollections.InsertOne(ctx, user)
        if insertErr != nil {
            msg := fmt.Sprint("user Item was not created ")
            c.JSON(http.StatusInternalServerError, gin.H{"error1": msg})
            return
        }
        defer cancel()
        c.JSON(http.StatusOK, gin.H{"insertionnumber": resultInsertionNumber})
    }
}

func Login() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        var user models.User
        var foundUser models.User
        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        err := userCollections.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
        defer cancel()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password not found"})
            return
        }
        passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
        defer cancel()
        if passwordIsValid != true {
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
        }
        if foundUser.Email == nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "user not found"})
            return
        }
        token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.UserType, &foundUser.UserId)
        helpers.UpdateAllTokens(token, refreshToken, foundUser.UserId)
        err = userCollections.FindOne(ctx, bson.M{"userid": foundUser.UserId}).Decode(&foundUser)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        c.JSON(http.StatusOK, foundUser)
    }
}

func GetUsers() gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
        recordsPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
        if err != nil || recordsPerPage < 1 {
            recordsPerPage = 10
        }
        page, _ := strconv.Atoi(c.Query("page"))
        if err != nil || page < 1 {
            page = 1
        }
        startIndex := (page - 1) * recordsPerPage
        startIndex, err = strconv.Atoi(c.Query("startIndex"))
        if err != nil {
            // handle error
        }
        matchStage := bson.D{{"$match", bson.D{}}}
        groupStage := bson.D{{"$group", bson.D{
            {"_id", "null"},
            {"total_count", bson.D{{"$sum", 1}}},
            {"data", bson.D{{"$push", "$$ROOT"}}},
        }}}
        projectStage := bson.D{
            {"$project", bson.D{
                {"_id", 0},
                {"total_count", 1},
                {"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordsPerPage}}}},
            }},
        }
        result, err := userCollections.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
        defer cancel()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing user items"})
        }
        var allUsers = []bson.M{}
        if err := result.All(ctx, &allUsers); err != nil {
            log.Fatal(err)
        }
        c.JSON(http.StatusOK, allUsers[0])
    }
}

func GetUser() gin.HandlerFunc {
    return func(c *gin.Context) {
        userId := c.Param("user_id")
        if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        var user models.User
        err := userCollections.FindOne(ctx, bson.M{"userid": userId}).Decode(&user)
        defer cancel()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, user)
    }
}

func Promote() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check if the user is an admin
        userType, exists := c.Get("usertype")
        if !exists || userType != "ADMIN" {
            c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
            return
        }

        // Get the user ID from the request
        userId := c.Param("user_id")
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        // Update the user type to ADMIN
        filter := bson.M{"userid": userId}
        update := bson.M{"$set": bson.M{"usertype": "ADMIN", "updatedat": time.Now()}}

        res, err := userCollections.UpdateOne(ctx, filter, update)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user type"})
            return
        }


        if res.MatchedCount == 0 {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        if res.ModifiedCount == 0 {
            c.JSON(http.StatusOK, gin.H{"message": "User type was already ADMIN"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin successfully"})
    }
}