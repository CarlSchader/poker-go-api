package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var size int64 = 2598960

func main() {
	mongodbURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("RANKS_COLLECTION_NAME")
	calcCollectionName := os.Getenv("CALCULATIONS_COLLECTION_NAME")
	batches, err := strconv.Atoi(os.Getenv("BATCHES"))

	if err != nil {
		log.Fatalln(errors.New("batches variable not set"))
	}

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

	// get collection and create handString index
	collection := client.Database(dbName).Collection(collectionName)
	calcCollection := client.Database(dbName).Collection(calcCollectionName)
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatalln(err)
	}

	// exit if the ranks collection already has at least all the rank documents
	if count >= size {
		fmt.Println("ranks collection already populated")
		os.Exit(0)
	}

	fmt.Println("creating indexes")
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"hand": 1,
		},
		Options: nil,
	})

	calcCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"hand":          1,
			"end_hand_size": 1,
		},
		Options: nil,
	})

	jsonBytes, err := ioutil.ReadFile(ranksFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	var ranksMap map[string]int
	err = json.Unmarshal(jsonBytes, &ranksMap)
	if err != nil {
		log.Fatalln(err)
	}

	var ranksObjects []interface{}
	for handString, rank := range ranksMap {
		ranksObjects = append(ranksObjects, bson.M{"hand": handString, "rank": rank})
	}

	fmt.Println("inserting data into db")
	entriesPerInsert := int(size) / batches
	extras := int(size) % batches
	for i := 0; i < batches; i++ {
		_, err := collection.InsertMany(context.Background(), ranksObjects[i*entriesPerInsert:(i+1)*entriesPerInsert])
		if err != nil {
			fmt.Printf("Error when writing %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("progress: %d of %d\n", i+1, batches)
	}

	if extras > 0 {
		_, err = collection.InsertMany(context.Background(), ranksObjects[batches*entriesPerInsert:batches*entriesPerInsert+extras])
		if err != nil {
			fmt.Printf("Error when writing extras %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("success")
}
