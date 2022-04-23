package database

import (
	"context"
	"fmt"
	"time"

	"github.com/carlschader/poker-go-api/application/algorithm"
	"github.com/carlschader/poker-go-api/application/simulate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Entry struct {
	Hand string `bson:"hand"`
	Rank int64  `bson:"rank"`
}

var client mongo.Client
var db mongo.Database
var cache mongo.Collection
var pockets mongo.Collection

func Connect(mongodbURI string, dbName string, cacheCollectionName string, pocketsCollectionName string) error {
	fmt.Printf("attempting to connect to %s\n", mongodbURI)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	newClient, err := mongo.NewClient(options.Client().ApplyURI(mongodbURI))
	client = *newClient
	if err != nil {
		return err
	}

	// connection loop
	for true {
		if err = client.Connect(ctx); err == nil {
			if err := client.Ping(context.TODO(), nil); err == nil {
				fmt.Println("db connected")
				break
			}
		}

		fmt.Println(err)
		fmt.Println("couldn't connect to mongodb, trying again")
		time.Sleep(3 * time.Second)
	}

	// defer client.Disconnect(ctx)
	db = *client.Database(dbName)
	cache = *db.Collection(cacheCollectionName)
	pockets = *db.Collection(pocketsCollectionName)

	return nil
}

func CollectionSize(collectionName string) (int64, error) {
	coll := db.Collection(collectionName)
	count, err := coll.CountDocuments(context.TODO(), bson.M{})
	return count, err
}

func LoadRanks(collectionName string) (map[string]int64, error) {
	coll := db.Collection(collectionName)
	cursor, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	results := []Entry{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	table := map[string]int64{}
	for _, entry := range results {
		table[entry.Hand] = entry.Rank
	}

	return table, nil
}

func CacheCheck(hand algorithm.Hand, endHandSize int, receiver interface{}) error {
	err := cache.FindOne(context.TODO(), bson.M{"hand": hand.Hash(), "end_hand_size": endHandSize}).Decode(receiver)
	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func CacheInsert(result *simulate.SimulationResult) error {
	_, err := cache.InsertOne(context.TODO(), *result)
	return err
}

func GetPocket(pocket algorithm.Hand, receiver interface{}) error {
	err := pockets.FindOne(context.TODO(), bson.M{"hand": pocket.Hash()}).Decode(receiver)
	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}
