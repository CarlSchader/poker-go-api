package algorithm

import (
	"sort"
	"strings"

	"gonum.org/v1/gonum/stat/combin"
)

var valMap = map[int]string{
	1: "A", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "10", 11: "J", 12: "Q", 13: "K",
}

var suitMap = map[int]string{
	1: "C", 2: "D", 3: "H", 4: "S",
}

type Card struct {
	Value int
	Suit  int
}

type Hand []Card
type Hand2 [2]Card
type Hand5 [5]Card
type Hand7 [7]Card

func (card Card) toString() string {
	return valMap[card.Value] + suitMap[card.Suit]
}

func (hand Hand) hash() string {
	size := len(hand)
	var stringSlice []string
	for i := 0; i < size; i++ {
		stringSlice = append(stringSlice, hand[i].toString())
	}
	sort.Strings(stringSlice)
	return strings.Join(stringSlice, "")
}

func (hand Hand5) Hash() string {
	var stringSlice []string
	for i := 0; i < 5; i++ {
		stringSlice = append(stringSlice, hand[i].toString())
	}
	sort.Strings(stringSlice)
	return strings.Join(stringSlice, "")
}

func HandToString(hand []Card) string {
	handString := ""
	length := len(hand)
	for i := 0; i < length; i++ {
		handString += hand[i].toString()
		if i < length-1 {
			handString += " "
		}
	}
	return handString
}

func MakeDeck(exclusionSet map[Card]bool) Hand {
	var deck Hand
	for suit := 1; suit <= 4; suit++ {
		for value := 1; value <= 13; value++ {
			card := Card{value, suit}
			if _, inSet := exclusionSet[card]; !inSet {
				deck = append(deck, card)
			}
		}
	}
	return deck
}

func CardCombinations(deck Hand, k int) chan Hand {
	n := len(deck)
	total := combin.Binomial(n, k)
	ch := make(chan Hand, total)
	// count := 0

	go func() {
		for indexComb := range indexCombinations(n, k) {
			// go func(indexComb []int) {
			comb := make(Hand, k)
			for i := 0; i < k; i++ {
				comb[i] = deck[indexComb[i]]
			}
			ch <- comb
			// count++
			// if count == total {
			// 	close(ch)
			// }
			// }(indexComb)
		}
		close(ch)
	}()

	return ch
}
