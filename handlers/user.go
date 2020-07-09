package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/mngibso/blog-api/db"
	"github.com/mngibso/blog-api/models"
)

const AuthUserKey = "authUsername"

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

type UserAPI struct {
	DB db.UserStorer
}

func NewUserAPI(s db.UserStorer) UserAPI {
	return UserAPI{db.NewUserStore()}
}

// CreateUser adds a new user to the database. Password is hashed before being stored in the DB.
// username must be unique
func (u UserAPI) CreateUser(ctx *gin.Context) {
	var user models.CreateUserInput
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: err.Error(),
		})
		return
	}
	count, err := u.DB.CountUsers(ctx, user.Username)
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
	insertedID, err := u.DB.InsertOne(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	r := models.NewUserFromUserInput(insertedID, user)
	ctx.JSON(http.StatusCreated, r)
}

// DeleteUser deletes the user and the user's posts
func (u UserAPI) DeleteUser(ctx *gin.Context) {
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
	err := u.DB.DeleteOne(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Message: "database error: " + err.Error(),
		})
		return
	}

	// delete user posts
	err = db.NewPostStore().DeleteMany(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Message: "database error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, models.ApiResponse{})
}

// GetUserByUserName returns the user object for the given username.
func (u UserAPI) GetUserByUserName(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: "username is required",
		})
		return
	}
	user, status, err := u.GetUser(username)
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
func (u UserAPI) GetUser(username string) (user models.User, status int, err error) {
	user, err = u.DB.FindOne(context.Background(), username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, http.StatusNotFound, err
		}
		return user, http.StatusInternalServerError, errors.New("database error: " + err.Error())
	}
	return user, http.StatusOK, nil
}

// getAllUsers returns an array of all users in the database
func (u UserAPI) getAllUsers() ([]models.User, error) {
	users, err := u.DB.Find(context.Background())
	if err != nil {
		return users, errors.New("database error: " + err.Error())
	}
	return users, nil
}

// GetUsers returns array of users
func (u UserAPI) GetUsers(ctx *gin.Context) {
	users, err := u.getAllUsers()
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

func (u UserAPI) UpdateUser(ctx *gin.Context) {
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

	// user can only update themselves
	if username != user.Username {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: "usernames must match",
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
	err = u.DB.FindOneAndReplace(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Message: "database error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, models.ApiResponse{})
}
