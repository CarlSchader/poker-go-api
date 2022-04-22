package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/carlschader/poker-go-api/application/algorithm"
	"github.com/carlschader/poker-go-api/application/simulate"
)

type PocketJson struct {
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

	// create hand stack
	fmt.Println("creating possible pockets")
	handStack := []algorithm.Hand{}
	for hand := range algorithm.CardCombinations(deck, 2) {
		handStack = append(handStack, hand)
	}

	// simulate pockets

	// fmt.Println("simulating pockets")
	// var pockets []*simulate.SimulationResult
	// count := 0
	// for hand := range algorithm.CardCombinations(deck, 2) {
	// 	result, err := simulate.SimulateHand(hand, 7, map[algorithm.Card]bool{}, ranksMap)
	// 	if err != nil {
	// 		log.Println(algorithm.HandToString(hand))
	// 		log.Println(err)
	// 		os.Exit(1)
	// 	} else {
	// 		pockets = append(pockets, result)
	// 		count++
	// 		fmt.Printf("%d of %d\n", count, TOTAL)
	// 	}
	// }

	var pockets []*simulate.SimulationResult
	count := 0
	cpus := runtime.NumCPU()
	fmt.Printf("simulating pockets on %d cores \n", cpus)
	runtime.GOMAXPROCS(cpus)
	var wg sync.WaitGroup
	wg.Add(cpus)
	for i := 0; i < cpus; i++ {
		go func() {
			defer wg.Done()
			for len(handStack) > 0 {
				// pop stack
				hand := handStack[len(handStack)-1]
				handStack = handStack[:len(handStack)-1]
				// calculate
				result, err := simulate.SimulateHand(hand, 7, map[algorithm.Card]bool{}, ranksMap)
				if err != nil {
					log.Println(hand.Hash())
					log.Println(err)
					os.Exit(1)
				} else {
					pockets = append(pockets, result)
					count++
					fmt.Printf("%d of %d\n", count, TOTAL)
				}
			}
		}()
	}
	wg.Wait()

	// sort pockets
	fmt.Println("sorting pockets")
	sort.SliceStable(pockets, func(i, j int) bool {
		return pockets[i].ExpectedRank < pockets[j].ExpectedRank
	})

	// write json
	fmt.Println("writing json")
	pocketsJsonMap := map[string]PocketJson{}
	for i, result := range pockets {
		pocketsJsonMap[result.Hand] = PocketJson{
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
