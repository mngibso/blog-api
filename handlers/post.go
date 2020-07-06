package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mngibso/blog-api/models"
)

var data = models.ApiResponse{
	Code:    0,
	Type_:   "Don't know",
	Message: "Hi therew",
}

func CreatePost(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data)
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
