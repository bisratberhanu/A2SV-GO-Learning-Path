package infrastructure

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SignedDetails struct{
	Email           string
	FirstName       string
	LastName        string
	Uid          	string
	UserType        string
	jwt.StandardClaims
	
}


var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName *string , lastName *string, userType *string, uid *string)(signedToken string, signedRefreshToken string, err error ){
	claims:= &SignedDetails{
		Email: email,
		FirstName: *firstName,
		LastName: *lastName,
		UserType: *userType,
		Uid: *uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(),
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
		fmt.Println(claims.ExpiresAt, time.Now().Local().Unix())
		msg = fmt.Sprintf("token has expired ")
		return
	}
	fmt.Println(msg)
	return claims,msg
} 

