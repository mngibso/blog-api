package handlers

import (
	"github.com/mngibso/blog-api/db"
	// "go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mngibso/blog-api/models"
)

const PostCollection = "posts"

var data = map[string]string{"hi": "there"}

func CreatePost(ctx *gin.Context) {
	var post models.CreatePostInput
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: err.Error(),
		})
		return
	}
	authUsername := ctx.MustGet(AuthUserKey).(string)
	post.Username = authUsername
	post.CreatedAt = time.Now().Unix()
	res, err := db.DB.Collection(PostCollection).InsertOne(ctx, post)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	r := models.NewPostFromPostInput(res.InsertedID, post)
	ctx.JSON(http.StatusCreated, r)
}

func DeletePost(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data)
}

func GetPost(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data)
}

func GetPostById(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data)
}

func UpdatePost(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data)
}
