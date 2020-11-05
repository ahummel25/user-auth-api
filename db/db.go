package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbFile string = "file://mnt/storage/db/user_auth.db?cache=shared&_auth&_auth_user=admin&_auth_pass=admin&_auth_crypt=sha1"

// GetDBConnection will return a mongo client connection.
func GetDBConnection() *mongo.Client {
	var (
		client   *mongo.Client
		err      error
		mongoURI = "mongodb+srv://user_auth_admin:rlippi7-yxyeEr@userauthmongoclusterdev.yohvj.mongodb.net/default?retryWrites=true&w=majority"
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI)); err != nil {
		log.Fatal(err)
	}

	return client

	/*defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("default")
	usersCollection := db.Collection("users")

	result, err := usersCollection.InsertOne(ctx, bson.D{
		{Key: "user_id", Value: uuid.New().String()},
		{Key: "email", Value: "ahummel25@gmail.com"},
		{Key: "username", Value: "ahummel25"},
	})

	fmt.Println(result)

	filter := bson.M{"username": "ahummel25"}
	var result bson.M
	err = usersCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result["username"])

	defer cur.Close(ctx)

	for cur.Next(ctx) {

		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result["user_id"])
		fmt.Println(result["username"])
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}*/
}
