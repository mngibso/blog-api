package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/mngibso/blog-api/db"
	"github.com/mngibso/blog-api/models"
)

const UserCollection = "users"

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(ctx *gin.Context) {
	var user models.CreateUserInput
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	count, err := db.DB.Collection(UserCollection).
		CountDocuments(ctx, bson.D{{"username", user.Username}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Message: "database error: " + err.Error(),
		})
		return
	}
	if count != 0 {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: "username already exists",
		})
		return
	}

	// store hashed password
	hashed, err := hashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Message: "hashing error: " + err.Error(),
		})
		return
	}
	user.Password = hashed
	res, err := db.DB.Collection(UserCollection).InsertOne(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	r := models.NewUserFromUserInput(res.InsertedID, user)
	ctx.JSON(http.StatusCreated, r)
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
