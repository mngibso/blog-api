package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/mngibso/blog-api/db"
	"github.com/mngibso/blog-api/models"
)

const UserCollection = "users"
const AuthUserKey = "authUsername"

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CreateUser adds a new user to the database
func CreateUser(ctx *gin.Context) {
	var user models.CreateUserInput
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: err.Error(),
		})
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

// DeleteUser deletes the user and the user's posts
func DeleteUser(ctx *gin.Context) {
	username := ctx.Param("username")
	authUsername := ctx.MustGet(AuthUserKey)

	if username == "" {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: "username is required",
		})
		return
	}
	if username != authUsername {
		// user can only delete themselves
		ctx.JSON(http.StatusForbidden, models.ApiResponse{
			Message: "invalid username",
		})
		return
	}
	q := bson.D{{"username", username}}
	_, err := db.DB.Collection(UserCollection).
		DeleteOne(ctx, q)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Message: "database error: " + err.Error(),
		})
		return
	}

	// Delete user posts
	_, err = db.DB.Collection(PostCollection).
		DeleteMany(ctx, bson.D{{"username", username}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Message: "database error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, models.ApiResponse{})
}

// GetUserByUserName returns the user object for the given username.
func GetUserByUserName(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: "username is required",
		})
		return
	}
	user, status, err := GetUser(username)
	if err != nil {
		ctx.JSON(status, models.ApiResponse{
			Message: err.Error(),
		})
		return
	}
	// don't send password
	user.Password = ""
	ctx.JSON(http.StatusOK, user)
}

// GetUser fetches the user with `username` from the database.  returns `status` as a http status code
func GetUser(username string) (user models.User, status int, err error) {
	q := bson.D{{"username", username}}
	err = db.DB.Collection(UserCollection).FindOne(context.Background(), q).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, http.StatusNotFound, err
		}
		return user, http.StatusInternalServerError, errors.New("database error: " + err.Error())
	}
	return user, http.StatusOK, nil
}

// getAllUsers returns an array of all users in the database
func getAllUsers() ([]models.User, error) {
	users := []models.User{}
	q := bson.D{}
	cur, err := db.DB.Collection(UserCollection).Find(context.Background(), q)
	if err != nil {
		return users, errors.New("database error: " + err.Error())
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			return users, errors.New("database error: " + err.Error())
		}
		// don't send password
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return users, errors.New("database error: " + err.Error())
	}
	return users, nil
}

// GetUsers returns array of users
func GetUsers(ctx *gin.Context) {
	users, err := getAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Message: err.Error(),
		})
		return
	}
	for _, v := range users {
		v.Password = ""
	}
	ctx.JSON(http.StatusOK, users)
}

func UpdateUser(ctx *gin.Context) {
	username := ctx.Param("username")
	authUsername := ctx.MustGet(AuthUserKey)
	if username == "" {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: "username is required",
		})
		return
	}
	if username != authUsername {
		// user can only update themselves
		ctx.JSON(http.StatusForbidden, models.ApiResponse{
			Message: "invalid username",
		})
		return
	}

	var user models.CreateUserInput
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: err.Error(),
		})
		return
	}

	if username != user.Username {
		// user can only update themselves
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: "usernames must match",
		})
		return
	}
	q := bson.D{{"username", username}}
	replacedUser := models.User{}
	// store hashed password
	hashed, err := hashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Message: "hashing error: " + err.Error(),
		})
		return
	}
	user.Password = hashed
	err = db.DB.Collection(UserCollection).
		FindOneAndReplace(ctx, q, user).
		Decode(&replacedUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Message: "database error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, models.ApiResponse{})
}
