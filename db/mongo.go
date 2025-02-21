package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

func Connect() (*mongo.Database, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(5000)*time.Second)
	var err error
	var uri = "mongodb://"
	if os.Getenv("TMATE_MONGO_USERNAME") != "" {
		uri += os.Getenv("TMATE_MONGO_USERNAME") + ":" + os.Getenv("TMATE_MONGO_PASSWORD") + "@"
	}
	uri += os.Getenv("TMATE_MONGO_HOST") + ":" + os.Getenv("TMATE_MONGO_PORT")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		fmt.Println("failed when trying to connect to '" + os.Getenv("TMATE_MONGO_HOST") + ":" + os.Getenv("TMATE_MONGO_PORT") + "' as '" + os.Getenv("TMATE_MONGO_USERNAME") + "'")
		fmt.Println(fmt.Errorf("unable to connect to mongo database"))
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("failed when trying to connect to '" + os.Getenv("TMATE_MONGO_HOST") + ":" + os.Getenv("TMATE_MONGO_PORT") + "' as '" + os.Getenv("TMATE_MONGO_USERNAME") + "'")
		fmt.Println(fmt.Errorf("unable to reach mongo database"))
	}

	return client.Database("tmate"), err
}
