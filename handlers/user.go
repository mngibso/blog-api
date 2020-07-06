package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUser(ctx *gin.Context) {
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
