package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func getDBClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("testing").Collection("numbers")
	res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	if err != nil {
		panic(err)
	}
	id := res.InsertedID
	fmt.Println(id)
	/*
		cur, err := collection.Find(ctx, bson.D{})
		if err != nil { log.Fatal(err) }
		defer cur.Close(ctx)
		for cur.Next(ctx) {
		   var result bson.M
		   err := cur.Decode(&result)
		   if err != nil { log.Fatal(err) }
		   // do something with result....
		}
		if err := cur.Err(); err != nil {
		  log.Fatal(err)
		}
		For methods that return a single item, a SingleResult instance is returned:

		var result struct {
		    Value float64
		}
		filter := bson.M{"name": "pi"}
		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
		err = collection.FindOne(ctx, filter).Decode(&result)
		if err != nil {
		    log.Fatal(err)
		}
	*/
}
