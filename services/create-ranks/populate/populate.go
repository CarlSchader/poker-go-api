package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongodbURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")
	ranksFilePath := os.Args[1]

	if len(os.Args) != 2 {
		log.Fatalln(errors.New("invalid number of command line arguments (should be 1)"))
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatalln(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Disconnect(ctx)

	// get collection and create cardString index
	collection := client.Database(dbName).Collection(collectionName)
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"card": 1,
		},
		Options: nil,
	})

	fmt.Println("loading data")
	jsonBytes, err := ioutil.ReadFile(ranksFilePath)
	if err != nil {
		log.Fatalln(errors.New("error reading ranks file"))
	}

	var ranksMap map[string]int
	err = json.Unmarshal(jsonBytes, &ranksMap)

	var ranksObjects []interface{}

	for cardString, rank := range ranksMap {
		ranksObjects = append(ranksObjects, bson.M{"card": cardString, "rank": rank})
	}

	fmt.Println("inserting data into db")
	insertCtx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	_, err = collection.InsertMany(insertCtx, ranksObjects)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("success")
	}
}
