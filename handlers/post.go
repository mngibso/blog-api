package handlers

import (
	"context"
	"fmt"
	"github.com/mngibso/blog-api/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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
	authUsername := ctx.MustGet(AuthUserKey)
	postID, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	fmt.Printf("%s %v", err, postID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: "invalid id: " + err.Error(),
		})
		return
	}

	q := bson.M{"_id": postID}
	found := models.Post{}
	err = db.DB.Collection(PostCollection).FindOne(context.Background(), q).Decode(&found)
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
	if authUsername != found.Username {
		ctx.JSON(http.StatusForbidden,
			models.ApiResponse{Message: "post not owned by user"})
		return
	}

	db.DB.Collection(PostCollection).FindOneAndDelete(ctx, q)
	ctx.JSON(http.StatusOK, models.ApiResponse{Message: "post deleted"})
}

func GetPostById(ctx *gin.Context) {
	postID := ctx.Param("id")
	if postID == "" {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: "post id is required",
		})
		return
	}
	q := bson.D{{"_id", postID}}
	found := models.Post{}
	err := db.DB.Collection(PostCollection).FindOne(ctx, q).Decode(&found)
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

func GetPost(ctx *gin.Context) {
	q := bson.D{{}}
	username := ctx.Query("username")
	if username != "" {
		q = bson.D{{"username", username}}
	}
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}})
	posts := []models.Post{}
	cur, err := db.DB.Collection(PostCollection).Find(ctx, q, opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			models.ApiResponse{Message: "database error: " + err.Error()})
		return
	}
	defer cur.Close(ctx)
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		var post models.Post
		err := cur.Decode(&post)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,
				models.ApiResponse{Message: "database error: " + err.Error()})
			return
		}
		// don't send password
		posts = append(posts, post)
	}
	if err := cur.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError,
			models.ApiResponse{Message: "database error: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, posts)
}

func UpdatePost(ctx *gin.Context) {
	authUsername := ctx.MustGet(AuthUserKey)
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: "invalid id: " + err.Error(),
		})
		return
	}

	var post models.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ApiResponse{
			Message: err.Error(),
		})
		return
	}
	if post.Username != authUsername {
		// user can only update their posts
		ctx.JSON(http.StatusForbidden, models.ApiResponse{
			Message: "invalid id or username",
		})
		return
	}

	q := bson.D{{"_id", id}}
	db.DB.Collection(PostCollection).FindOneAndReplace(ctx, q, post)
	ctx.JSON(http.StatusOK, models.ApiResponse{Message: "post updated"})
}
