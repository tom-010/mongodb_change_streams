package main

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func iterateChangeStream(routineCtx context.Context, waitGroup sync.WaitGroup, stream *mongo.ChangeStream) {
	defer stream.Close(routineCtx)
	defer waitGroup.Done()
	for stream.Next(routineCtx) {
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			panic(err)
		}
		j, _ := bson.MarshalExtJSON(data["fullDocument"], true, false)

		fmt.Printf("%s\n", j)
	}
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	database := client.Database("quickstart")
	episodesCollection := database.Collection("episodes")

	var waitGroup sync.WaitGroup

	episodesStream, err := episodesCollection.Watch(context.TODO(), mongo.Pipeline{})
	if err != nil {
		panic(err)
	}
	waitGroup.Add(1)
	routineCtx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	go iterateChangeStream(routineCtx, waitGroup, episodesStream)

	waitGroup.Wait()
}
