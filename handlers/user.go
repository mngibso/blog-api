package handlers

import (
	"fmt"
	"github.com/mngibso/blog-api/db"
	"github.com/mngibso/blog-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"

	"github.com/gin-gonic/gin"
)

const UserCollection = "users"

func CreateUser(ctx *gin.Context) {
	var existing, user models.CreateUserInput
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// find one
	if err := db.DB.Collection(UserCollection).
		FindOne(ctx, bson.D{{"username", user.Username}}).
		Decode(&existing); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing.Username != "" {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "Username already exists",
		})
		return
	}

	res, err := db.DB.Collection(UserCollection).InsertOne(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(res)
	ctx.JSON(http.StatusOK, data)
}

func DeleteUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data)
}

func GetUserByUserName(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data)
}

func GetUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data)
}

func UserLogin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data)
}

func UserLogout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data)
}

func UpdateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data)
}
