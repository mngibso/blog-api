package routes

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/mngibso/blog-api/handlers"
	"github.com/mngibso/blog-api/models"
)

// basicAuth returns middleware that checks Basic Authentication
func basicAuth() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		auth := strings.SplitN(ctx.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			ctx.JSON(http.StatusUnauthorized, models.ApiResponse{
				Message: "authentication required",
			})
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !authenticateUser(pair[0], pair[1]) {
			ctx.Header("WWW-Authenticate", "Authentication Required")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// set user in context, can be obtained later with ctx.MustGet(handler.AuthUserKey).
		ctx.Set(handlers.AuthUserKey, pair[0])
		ctx.Next()
	}
}

// authenticateUser checks credentials by comparing password in header with password in db
func authenticateUser(username string, password string) bool {
	user, _, err := handlers.GetUser(username)
	if err != nil {
		log.Printf("unable to authenticate user, error: %s", err.Error())
		return false
	}
	return checkPasswordHash(password, user.Password)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
