package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/carlschader/poker-go-api/application/algorithm"
	"github.com/carlschader/poker-go-api/application/database"
	"github.com/carlschader/poker-go-api/application/simulate"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type PocketResponse struct {
	Hand             string                     `json:"hand" bson:"hand"`
	Rank             int                        `json:"rank" bson:"rank"`
	MaxRank          int                        `json:"max_rank" bson:"max_rank"`
	Percentile       float32                    `json:"percentile" bson:"percentile"`
	SimulationResult *simulate.SimulationResult `json:"simulation_result" bson:"simulation_result"`
}

var ranks map[string]int64

func queryToHand(handString string) algorithm.Hand {
	hand := algorithm.Hand{}
	for _, cardString := range strings.Split(handString, "-") {
		if len(cardString) == 2 {
			hand = append(hand, algorithm.Card{
				Value: algorithm.InvertedValMap[cardString[:1]],
				Suit:  algorithm.InvertedSuitMap[cardString[1:]],
			})
		}
	}
	return hand
}

func main() {
	port := os.Getenv("PORT")
	mongodbURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")
	ranksCollectionName := os.Getenv("RANKS_COLLECTION_NAME")
	cacheCollectionName := os.Getenv("CACHE_COLLECTION_NAME")
	pocketsCollectionName := os.Getenv("POCKETS_COLLECTION_NAME")
	var calculationTimeout time.Duration

	if secondsInt, err := strconv.Atoi(os.Getenv("CALCULATION_TIMEOUT")); err != nil {
		log.Fatalln(err)
	} else {
		calculationTimeout = time.Second * time.Duration(secondsInt)
		log.Printf("timeout allowed: %v\n", calculationTimeout)
	}

	if err := database.Connect(mongodbURI, dbName, cacheCollectionName, pocketsCollectionName); err != nil {
		log.Fatalln(err)
	}

	app := fiber.New()

	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		startTime := time.Now()
		count, err := database.CollectionSize(ranksCollectionName)
		if err != nil {
			log.Println(err)
			if err == mongo.ErrClientDisconnected {
				go func() {
					if err := database.Connect(mongodbURI, dbName, cacheCollectionName, pocketsCollectionName); err != nil {
						log.Fatalln(err)
					}
				}()
			}
			return err
		}

		if count < 2598960 {
			return c.JSON(map[string]string{"message": "ranks collection not populated, please wait"})
		} else if len(ranks) < 2598960 {
			go func() {
				var err error
				ranks, err = database.LoadRanks(ranksCollectionName)
				if err != nil {
					log.Fatal(err)
				}
			}()
			return c.JSON(map[string]string{"message": "ranks not cached in memory, please wait"})
		}

		hand := queryToHand(strings.ToUpper(c.Query("hand")))
		shared := queryToHand(strings.ToUpper(c.Query("shared")))
		endHandSize, err := strconv.Atoi(c.Query("count"))
		if err != nil {
			log.Println(err)
			return err
		}

		seenCards := map[algorithm.Card]bool{}
		for _, card := range hand {
			if seen, ok := seenCards[card]; ok && seen {
				return c.JSON("invalid hand")
			}
			seenCards[card] = true
		}
		for _, card := range shared {
			if seen, ok := seenCards[card]; ok && seen {
				return c.JSON("invalid hand")
			}
			seenCards[card] = true
		}

		if len(hand) != 2 {
			return errors.New("invalid hand size. Must be two")
		}

		if len(shared) == 0 {
			var result PocketResponse
			err := database.GetPocket(hand, &result)
			if err != nil {
				log.Println(err)
				return err
			}
			result.MaxRank = 1326
			result.Percentile = (float32(result.Rank) / float32(result.MaxRank)) * 100.0
			log.Printf("calculation time: %v", time.Since(startTime))
			return c.JSON(result)
		} else {
			var joinedHand algorithm.Hand
			for _, card := range hand {
				joinedHand = append(joinedHand, card)
			}
			for _, card := range shared {
				joinedHand = append(joinedHand, card)
			}

			var result *simulate.SimulationResult
			err := database.CacheCheck(joinedHand, endHandSize, &result)
			if result == nil && err == nil {
				log.Println("calculation not cached yet")
				result, err = simulate.SimulateHand(joinedHand, endHandSize, map[algorithm.Card]bool{}, ranks)
				if err != nil {
					log.Println(err)
					return err
				}

				err = database.CacheInsert(result)
				if err != nil {
					log.Println(err)
					return err
				}
			} else if err != nil {
				log.Println(err)
				return err
			}
			log.Printf("calculation time: %v", time.Since(startTime))
			return c.JSON(result)
		}
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Listen(fmt.Sprintf(":%s", port))
}
