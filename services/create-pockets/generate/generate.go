package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"

	"github.com/carlschader/poker-go-api/application/algorithm"
	"github.com/carlschader/poker-go-api/application/simulate"
)

type pocketJson struct {
	Rank             int                        `json:"rank" bson:"rank"`
	SimulationResult *simulate.SimulationResult `json:"simulation_result" bson:"simulation_result"`
}

func main() {
	start := time.Now()

	TOTAL := 1326

	if len(os.Args) < 3 {
		fmt.Println(errors.New("must specify a file path for ranks file and a destination file as args"))
		os.Exit(1)
	}
	ranksFilePath := os.Args[1]

	deck := algorithm.MakeDeck(map[algorithm.Card]bool{})

	// load ranks
	fmt.Println("loading ranks table")
	jsonBytes, err := ioutil.ReadFile(ranksFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	var ranksMap map[string]int64
	err = json.Unmarshal(jsonBytes, &ranksMap)
	if err != nil {
		log.Fatalln(err)
	}

	// simulate pockets
	fmt.Println("simulating pockets")
	var pockets []*simulate.SimulationResult
	count := 0
	// var wg sync.WaitGroup
	// for hand := range algorithm.CardCombinations(deck, 2) {
	// 	wg.Add(1)
	// 	go func(hand algorithm.Hand, count int, wg sync.WaitGroup) {
	// 		defer wg.Done()
	// 		result, err := simulate.SimulateHand(hand, 7, map[algorithm.Card]bool{}, ranksMap)
	// 		if err != nil {
	// 			log.Println(algorithm.HandToString(hand))
	// 			log.Println(err)
	// 			os.Exit(1)
	// 		} else {
	// 			pockets = append(pockets, result)
	// 			count++
	// 			fmt.Printf("%d of %d\n", count, TOTAL)
	// 		}
	// 	}(hand, count, wg)
	// }
	// wg.Wait()

	for hand := range algorithm.CardCombinations(deck, 2) {
		result, err := simulate.SimulateHand(hand, 7, map[algorithm.Card]bool{}, ranksMap)
		if err != nil {
			log.Println(algorithm.HandToString(hand))
			log.Println(err)
			os.Exit(1)
		} else {
			pockets = append(pockets, result)
			count++
			fmt.Printf("%d of %d\n", count, TOTAL)
		}
	}

	// sort pockets
	fmt.Println("sorting pockets")
	sort.SliceStable(pockets, func(i, j int) bool {
		return pockets[i].ExpectedRank < pockets[j].ExpectedRank
	})

	// write json
	fmt.Println("writing json")
	pocketsJsonMap := map[string]pocketJson{}
	for i, result := range pockets {
		pocketsJsonMap[result.Hand] = pocketJson{
			i + 1,
			result,
		}
	}

	jsonString, err := json.MarshalIndent(pocketsJsonMap, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(os.Args[2], jsonString, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("done\n")
	fmt.Println(time.Since(start))
}
