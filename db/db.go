package db

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	InsertOne(ctx context.Context, user models.CreateUserInput) (interface{}, error)
	DeleteOne(ctx context.Context, username string) error
	FindOne(ctx context.Context, username string) (models.User, error)
	FindOneAndReplace(ctx context.Context, user models.CreateUserInput) error
	Find(ctx context.Context) ([]models.User, error)
}

type PostStorer interface {
	DeleteMany(ctx context.Context, username string) error
	InsertOne(ctx context.Context, post models.CreatePostInput) (interface{}, error)
	Find(ctx context.Context, username string) ([]models.Post, error)
	FindOne(ctx context.Context, postID string) (models.Post, error)
	FindOneAndDelete(ctx context.Context, postID string) error
	FindOneAndReplace(ctx context.Context, post models.Post) error
}

type PostStore struct {
	Collection string
}

type UserStore struct {
	Collection string
}

func NewUserStore() UserStorer {
	return &UserStore{UserCollection}
}

func NewPostStore() PostStorer {
	return &PostStore{PostCollection}
}

func (s UserStore) CountUsers(ctx context.Context, username string) (int64, error) {
	return DB.Collection(s.Collection).
		CountDocuments(ctx, bson.D{{"username", username}})
}

func (s UserStore) InsertOne(ctx context.Context, user models.CreateUserInput) (interface{}, error) {
	res, err := DB.Collection(s.Collection).InsertOne(ctx, user)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID, nil
}

func (s UserStore) DeleteOne(ctx context.Context, username string) error {
	q := bson.D{{"username", username}}
	_, err := DB.Collection(s.Collection).DeleteOne(ctx, q)
	return err
}

func (s UserStore) FindOne(ctx context.Context, username string) (models.User, error) {
	user := models.User{}
	q := bson.D{{"username", username}}
	err := DB.Collection(s.Collection).FindOne(ctx, q).Decode(&user)
	return user, err
}

func (s UserStore) FindOneAndReplace(ctx context.Context, user models.CreateUserInput) error {
	replacedUser := models.User{}
	q := bson.D{{"username", user.Username}}
	return DB.Collection(s.Collection).
		FindOneAndReplace(ctx, q, user).
		Decode(&replacedUser)
}

func (s UserStore) Find(ctx context.Context) (users []models.User, err error) {
	q := bson.D{}
	cur, err := DB.Collection(s.Collection).Find(ctx, q)
	if err != nil {
		return users, errors.New("database error: " + err.Error())
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		// To decode into a struct, use cursor.Decode()
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			return users, errors.New("database error: " + err.Error())
		}
		// don't send password
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return users, err
	}
	return users, nil
}

// Initialize mongodb connection
func InitializeMongoDB() {
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

func (p PostStore) DeleteMany(ctx context.Context, username string) error {
	_, err := DB.Collection(p.Collection).
		DeleteMany(ctx, bson.D{{"username", username}})
	return err
}

func (p PostStore) InsertOne(ctx context.Context, post models.CreatePostInput) (interface{}, error) {
	res, err := DB.Collection(p.Collection).InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, err
}

func (p PostStore) Find(ctx context.Context, username string) ([]models.Post, error) {
	// , opts options.FindOptions
	q := bson.D{{}}
	if username != "" {
		q = bson.D{{"username", username}}
	}
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}})
	posts := []models.Post{}
	cur, err := DB.Collection(p.Collection).Find(ctx, q, opts)
	if err != nil {
		return posts, err
	}
	defer cur.Close(ctx)
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		var post models.Post
		err := cur.Decode(&post)
		if err != nil {
			return posts, err
		}
		// don't send password
		posts = append(posts, post)
	}
	err = cur.Err()
	return posts, err
}

func (p PostStore) FindOne(ctx context.Context, pID string) (post models.Post, err error) {
	postID, err := primitive.ObjectIDFromHex(pID)
	if err != nil {
		return post, err
	}
	q := bson.M{"_id": postID}
	err = DB.Collection(p.Collection).FindOne(ctx, q).Decode(&post)
	return post, err
}

func (p PostStore) FindOneAndDelete(ctx context.Context, postID string) error {
	post := models.Post{}
	pID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	q := bson.M{"_id": pID}
	return DB.Collection(p.Collection).FindOneAndDelete(ctx, q).Decode(&post)
}

func (p PostStore) FindOneAndReplace(ctx context.Context, post models.Post) error {
	q := bson.M{"_id": post.ID}
	replacePost := models.Post{}
	return DB.Collection(p.Collection).FindOneAndReplace(ctx, q, post).Decode(&replacePost)
}
