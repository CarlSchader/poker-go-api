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

	"github.com/carlschader/poker-go-api/application/simulate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PocketJson struct {
	Rank             int                        `json:"rank" bson:"rank"`
	SimulationResult *simulate.SimulationResult `json:"simulation_result" bson:"simulation_result"`
}

var size int64 = 1326

func main() {
	mongodbURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("POCKETS_COLLECTION_NAME")
	batches := 1

	pocketsFilePath := os.Args[1]

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

	// delete documents that are already in the collection
	collection.DeleteMany(ctx, bson.M{})

	fmt.Println("creating indexes")
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"hand": 1,
		},
		Options: nil,
	})

	jsonBytes, err := ioutil.ReadFile(pocketsFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	var pocketsMap map[string]PocketJson
	err = json.Unmarshal(jsonBytes, &pocketsMap)
	if err != nil {
		log.Fatalln(err)
	}

	var pocketsObjects []interface{}
	for handString, entry := range pocketsMap {
		pocketsObjects = append(pocketsObjects, bson.M{"hand": handString, "rank": entry.Rank, "simulation_result": entry.SimulationResult})
	}

	fmt.Println("inserting data into db")
	entriesPerInsert := int(size) / batches
	extras := int(size) % batches
	for i := 0; i < batches; i++ {
		_, err := collection.InsertMany(context.Background(), pocketsObjects[i*entriesPerInsert:(i+1)*entriesPerInsert])
		if err != nil {
			fmt.Printf("Error when writing %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("progress: %d of %d\n", i+1, batches)
	}

	if extras > 0 {
		_, err = collection.InsertMany(context.Background(), pocketsObjects[batches*entriesPerInsert:batches*entriesPerInsert+extras])
		if err != nil {
			fmt.Printf("Error when writing extras %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("success")
}
