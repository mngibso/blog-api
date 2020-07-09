package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mngibso/blog-api/db"
	"github.com/mngibso/blog-api/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreatePost creates a blog post in the database
func CreatePost(ctx *gin.Context) {
	coll := db.NewPostStore()
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

	insertedID, err := coll.InsertOne(ctx, post)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	r := models.NewPostFromPostInput(insertedID, post)
	ctx.JSON(http.StatusCreated, r)
}

// DeletePost deletes a post from the db
func DeletePost(ctx *gin.Context) {
	coll := db.NewPostStore()
	authUsername := ctx.MustGet(AuthUserKey)
	postID := ctx.Param("id")

	post, err := coll.FindOne(ctx, postID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, models.ApiResponse{Message: "not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError,
			models.ApiResponse{Message: "database error: " + err.Error()})
		return
	}

	// user can only delete their own post
	if authUsername != post.Username {
		ctx.JSON(http.StatusForbidden,
			models.ApiResponse{Message: "post not owned by user"})
		return
	}

	coll.FindOneAndDelete(ctx, postID)
	ctx.JSON(http.StatusOK, models.ApiResponse{Message: "post deleted"})
}

func GetPostById(ctx *gin.Context) {
	coll := db.NewPostStore()
	postID := ctx.Param("id")
	if postID == "" {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: "post id is required",
		})
		return
	}
	found, err := coll.FindOne(ctx, postID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, models.ApiResponse{Message: "not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError,
			models.ApiResponse{Message: "database error: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, found)
}

// GetPost returns all posts in the DB. `username` is an optional
// query string parameter to filter results for a user.  Sorts by createdAt, descending
func GetPost(ctx *gin.Context) {
	coll := db.NewPostStore()
	username := ctx.Query("username")
	posts, err := coll.Find(ctx, username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			models.ApiResponse{Message: "database error: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, posts)
}

// UpdatePost replaces an existing post with the request body
func UpdatePost(ctx *gin.Context) {
	coll := db.NewPostStore()
	authUsername := ctx.MustGet(AuthUserKey)
	var post models.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: err.Error(),
		})
		return
	}
	// user can only update their posts
	if post.Username != authUsername {
		ctx.JSON(http.StatusForbidden, models.ApiResponse{
			Message: "invalid id or username",
		})
		return
	}
	if err := coll.FindOneAndReplace(ctx, post); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Message: "database error: " + err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, models.ApiResponse{Message: "post updated"})
}
