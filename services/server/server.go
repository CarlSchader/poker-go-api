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

var ranks map[string]int64

func main() {
	port := os.Getenv("PORT")
	mongodbURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")
	ranksCollectionName := os.Getenv("RANKS_COLLECTION_NAME")
	cacheCollectionName := os.Getenv("CACHE_COLLECTION_NAME")
	var calculationTimeout time.Duration

	if secondsInt, err := strconv.Atoi(os.Getenv("CALCULATION_TIMEOUT")); err != nil {
		log.Fatalln(err)
	} else {
		calculationTimeout = time.Second * time.Duration(secondsInt)
		log.Printf("timeout allowed: %v\n", calculationTimeout)
	}

	if err := database.Connect(mongodbURI, dbName, cacheCollectionName); err != nil {
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
					if err := database.Connect(mongodbURI, dbName, cacheCollectionName); err != nil {
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

		if len(hand) != 2 {
			return errors.New("invalid hand size. Must be two")
		}

		if len(shared) == 0 {

		}

		result, err := database.CacheCheck(hand, endHandSize)

		if result == nil && err == nil {
			log.Println("calculation not cached yet")
			result, err = simulate.SimulateHand(hand, endHandSize, map[algorithm.Card]bool{}, ranks)
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
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Listen(fmt.Sprintf(":%s", port))
}
