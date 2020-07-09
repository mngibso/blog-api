package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Database

const DBNAME = "blog"

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
