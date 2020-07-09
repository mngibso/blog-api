package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/mngibso/blog-api/models"
)

var DB *mongo.Database

const DBNAME = "blog"
const PostCollection = "posts"
const UserCollection = "users"

type UserStorer interface {
	CountUsers(ctx context.Context, username string) (int64, error)
	InsertOne(ctx context.Context, user models.CreateUserInput) (string, error)
	DeleteOne(ctx context.Context, q interface{}) error
	DeleteMany(ctx context.Context, q interface{}) error
	FindOne(ctx context.Context, q interface{}) (models.User, error)
	FindOneAndReplace(ctx context.Context, q interface{}, user models.CreateUserInput) error
	Find(ctx context.Context, query interface{}) ([]models.User, error)
}

type PostStorer interface {
	InsertOne(ctx context.Context, user models.CreatePostInput) (string, error)
	Find(ctx context.Context, query interface{}, opts options.FindOptions) ([]models.Post, error)
	FindOne(ctx context.Context, q interface{}) (models.Post, error)
	FindOneAndDelete(ctx context.Context, q interface{})
	FindOneAndReplace(ctx context.Context, q interface{}, post models.Post)
}

// Initialize mongodb connection
func init() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("`MONGODB_URI must be set in the env")
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	DB = client.Database(DBNAME)
}
