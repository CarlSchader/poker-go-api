package algorithm

import (
	"fmt"
	"testing"
)

func TestMakeDeck(t *testing.T) {
	deck := MakeDeck(map[Card]bool{})
	fmt.Printf("%s\n deck size: %d\n", HandToString(deck), len(deck))

	deck = MakeDeck(map[Card]bool{
		{2, 1}: true,
		{2, 2}: true,
		{2, 3}: true,
		{2, 4}: true,
	})
	fmt.Printf("%s\n deck size: %d\n", HandToString(deck), len(deck))
}

func TestCardCombinations(t *testing.T) {
	deck := MakeDeck(map[Card]bool{
		{2, 1}: true,
		{2, 2}: true,
		{2, 3}: true,
		{2, 4}: true,
	})
	i := 0
	for hand := range CardCombinations(deck, 2) {
		i++
		fmt.Printf("hand%d: %v\n", i, hand)
	}
}

func TestCardCompare(t *testing.T) {
	hand := Hand{Card{1, 2}, Card{1, 3}, Card{1, 4}, Card{2, 2}, Card{2, 3}}
	for comb := range CardCombinations(Hand{Card{1, 1}, Card{2, 4}, Card{7, 3}, Card{8, 3}, Card{9, 3}, Card{10, 3}, Card{11, 3}, Card{13, 2}}, 5) {
		fmt.Printf("hand: %v, comb %v\n", hand, comb)
		if HandLessThan(hand, comb) {
			fmt.Printf("\twinner %v\n", hand)
		} else {
			fmt.Printf("\twinner %v\n", comb)
		}
	}
}
