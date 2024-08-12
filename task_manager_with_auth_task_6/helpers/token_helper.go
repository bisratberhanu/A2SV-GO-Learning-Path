package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"task_manager_with_auth_task_6/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct{
	Email           string
	FirstName       string
	LastName        string
	Uid          	string
	UserType        string
	jwt.StandardClaims
	
}

var userCollections *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName *string , lastName *string, userType *string, uid *string)(signedToken string, signedRefreshToken string, err error ){
	claims:= &SignedDetails{
		Email: email,
		FirstName: *firstName,
		LastName: *lastName,
		UserType: *userType,
		Uid: *uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshclaims:= &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token,err:= jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString([]byte(SECRET_KEY))
	refreshToken,err := jwt.NewWithClaims(jwt.SigningMethodHS256,refreshclaims).SignedString([]byte(SECRET_KEY))

	if err!=nil{
		log.Panic()
		return
	}

	return token,refreshToken,err
}


func UpdateAllTokens(token string, refreshToken string, userId string){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
	var updateObject primitive.D
	updatedAt,_ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObject = append(updateObject, bson.E{Key: "updatedat",  Value: updatedAt })
	upsert:= true
	filter:= bson.M{"userid": userId}
	opt:= options.UpdateOptions{Upsert: &upsert,}
_,err := userCollections.UpdateOne(ctx, filter,bson.D{{"$set", updateObject}},&opt)
	defer cancel()
	if err!=nil{
		log.Panic(err)
		return
	}
	
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string){
	token, err := jwt.ParseWithClaims(
		signedToken,&SignedDetails{}, func(token *jwt.Token)(interface{},error){
			return []byte(SECRET_KEY),nil
		},
	)
	if err!=nil{
		msg = err.Error()
	}
	claims,ok := token.Claims.(*SignedDetails)
	if !ok{
		msg = fmt.Sprintf("token is not valid")
		msg = err.Error()
		return

	}
	if claims.ExpiresAt < time.Now().Local().Unix(){
		msg = fmt.Sprintf("token has expired ")
		return
	}
	return claims,msg
} 
