package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mngibso/blog-api/db"
	"github.com/mngibso/blog-api/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostAPI struct {
	DB db.PostStorer
}

func NewPostAPI(p db.PostStorer) PostAPI {
	return PostAPI{p}
}

// CreatePost creates a blog post in the database
func (p PostAPI) CreatePost(ctx *gin.Context) {
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

	insertedID, err := p.DB.InsertOne(ctx, post)
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
func (p PostAPI) DeletePost(ctx *gin.Context) {
	authUsername := ctx.MustGet(AuthUserKey)
	postID := ctx.Param("id")

	post, err := p.DB.FindOne(ctx, postID)
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

	p.DB.FindOneAndDelete(ctx, postID)
	ctx.JSON(http.StatusOK, models.ApiResponse{Message: "post deleted"})
}

func (p PostAPI) GetPostById(ctx *gin.Context) {
	postID := ctx.Param("id")
	if postID == "" {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: "post id is required",
		})
		return
	}
	found, err := p.DB.FindOne(ctx, postID)
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
func (p PostAPI) GetPost(ctx *gin.Context) {
	username := ctx.Query("username")
	posts, err := p.DB.Find(ctx, username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			models.ApiResponse{Message: "database error: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, posts)
}

// UpdatePost replaces an existing post with the request body
func (p PostAPI) UpdatePost(ctx *gin.Context) {
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
	if err := p.DB.FindOneAndReplace(ctx, post); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ApiResponse{
			Message: "database error: " + err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, models.ApiResponse{Message: "post updated"})
}
