package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/carlschader/poker-go-api/server/algorithm"
)

func main() {
	start := time.Now()
	deck := algorithm.MakeDeck(map[algorithm.Card]bool{})
	unsortedMap := [10][]algorithm.Hand5{}
	fmt.Println("generating")
	for hand := range algorithm.CardCombinations(deck, 5) {
		if best := algorithm.GetRoyalFlush(hand); best[0].Value != 0 {
			unsortedMap[0] = append(unsortedMap[0], best)
		} else if best := algorithm.GetStraightFlush(hand); best[0].Value != 0 {
			unsortedMap[1] = append(unsortedMap[1], best)
		} else if best := algorithm.GetFourOfAKind(hand); best[0].Value != 0 {
			unsortedMap[2] = append(unsortedMap[2], best)
		} else if best := algorithm.GetFullHouse(hand); best[0].Value != 0 {
			unsortedMap[3] = append(unsortedMap[3], best)
		} else if best := algorithm.GetFlush(hand); best[0].Value != 0 {
			unsortedMap[4] = append(unsortedMap[4], best)
		} else if best := algorithm.GetStraight(hand); best[0].Value != 0 {
			unsortedMap[5] = append(unsortedMap[5], best)
		} else if best := algorithm.GetThreeOfAKind(hand); best[0].Value != 0 {
			unsortedMap[6] = append(unsortedMap[6], best)
		} else if best := algorithm.GetTwoPair(hand); best[0].Value != 0 {
			unsortedMap[7] = append(unsortedMap[7], best)
		} else if best := algorithm.GetPair(hand); best[0].Value != 0 {
			unsortedMap[8] = append(unsortedMap[8], best)
		} else {
			best := algorithm.GetFiveHighestCards(hand)
			unsortedMap[9] = append(unsortedMap[9], best)
		}
	}
	total := 0
	sort.SliceStable(unsortedMap[0], func(i, j int) bool {
		return algorithm.RoyalFlushLessThan(unsortedMap[0][i], unsortedMap[0][j])
	})
	fmt.Printf("sorting %d\n", 1)
	sort.SliceStable(unsortedMap[1], func(i, j int) bool {
		return algorithm.StraightFlushLessThan(unsortedMap[1][i], unsortedMap[1][j])
	})
	fmt.Printf("sorting %d\n", 2)
	sort.SliceStable(unsortedMap[2], func(i, j int) bool {
		return algorithm.FourOfAKindLessThan(unsortedMap[2][i], unsortedMap[2][j])
	})
	fmt.Printf("sorting %d\n", 3)
	sort.SliceStable(unsortedMap[3], func(i, j int) bool {
		return algorithm.FullHouseLessThan(unsortedMap[3][i], unsortedMap[3][j])
	})
	fmt.Printf("sorting %d\n", 4)
	sort.SliceStable(unsortedMap[4], func(i, j int) bool {
		return algorithm.FlushLessThan(unsortedMap[4][i], unsortedMap[4][j])
	})
	fmt.Printf("sorting %d\n", 5)
	sort.SliceStable(unsortedMap[5], func(i, j int) bool {
		return algorithm.StraightLessThan(unsortedMap[5][i], unsortedMap[5][j])
	})
	fmt.Printf("sorting %d\n", 6)
	sort.SliceStable(unsortedMap[6], func(i, j int) bool {
		return algorithm.ThreeOfAKindLessThan(unsortedMap[6][i], unsortedMap[6][j])
	})
	fmt.Printf("sorting %d\n", 7)
	sort.SliceStable(unsortedMap[7], func(i, j int) bool {
		return algorithm.TwoPairLessThan(unsortedMap[7][i], unsortedMap[7][j])
	})
	fmt.Printf("sorting %d\n", 8)
	sort.SliceStable(unsortedMap[8], func(i, j int) bool {
		return algorithm.PairLessThan(unsortedMap[8][i], unsortedMap[8][j])
	})
	fmt.Printf("sorting %d\n", 9)
	sort.SliceStable(unsortedMap[9], func(i, j int) bool {
		return algorithm.HighCardLessThan(unsortedMap[9][i], unsortedMap[9][j])
	})

	for i := 0; i < 10; i++ {
		total += len(unsortedMap[i])
	}
	fmt.Printf("sorting %d\n", 10)
	fmt.Println("writing")

	rank := 0
	ranked := map[string]int{}
	for i := 9; i >= 0; i-- {
		for _, hand := range unsortedMap[i] {
			rank++
			ranked[hand.Hash()] = rank
		}
	}
	jsonString, err := json.MarshalIndent(ranked, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(os.Args[1], jsonString, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("done\n")
	fmt.Println(time.Since(start))
	fmt.Printf("total: %d\n", total)
}
