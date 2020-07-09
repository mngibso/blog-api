package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mngibso/blog-api/mocks"
	"github.com/mngibso/blog-api/models"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPostAPI(t *testing.T) {
	Convey("Given a BTree", t, func() {
		postStore := mocks.PostStorer{}
		postAPI := NewPostAPI(&postStore)
		userName1 := "allen"
		post1 := models.CreatePostInput{
			Title: "Blog Post 1",
			Body:  "Body For Blog Post 1",
		}
		Convey("POST'ing a blog post should create the blog post", func() {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Set(AuthUserKey, userName1)
			p1, _ := json.Marshal(post1)
			// Create the post request with post body
			request, _ := http.NewRequest("POST", "/", bytes.NewBuffer(p1))
			ctx.Request = request
			insertedID, _ := primitive.ObjectIDFromHex("1")
			postStore.On("InsertOne", ctx, mock.Anything).Return(insertedID, nil)

			// call the handler
			postAPI.CreatePost(ctx)
			So(ctx.IsAborted(), ShouldBeFalse)
			So(len(ctx.Errors), ShouldEqual, 0)
			So(ctx.Writer.Status(), ShouldEqual, http.StatusCreated)

			// test return from request
			post := models.Post{}
			json.Unmarshal(recorder.Body.Bytes(), &post)
			So(post.Username, ShouldEqual, userName1)
			So(post.Body, ShouldEqual, post1.Body)
			So(post.Title, ShouldEqual, post1.Title)
			So(post.CreatedAt, ShouldNotEqual, 0)
		})
	})
}
