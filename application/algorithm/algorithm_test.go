package algorithm

import (
	"fmt"
	"testing"
)

func TestGetRoyalFlush(t *testing.T) {
	deck := MakeDeck(map[Card]bool{
		{1, 4}: true, {13, 4}: true, {12, 4}: true, {11, 4}: true, {10, 4}: true,
	})
	hand := GetRoyalFlush(deck)
	fmt.Printf("royal flush: %v\n", hand)
}

func TestGetStraightsAndFlushes(t *testing.T) {
	deck := MakeDeck(map[Card]bool{
		{1, 4}: true, {13, 4}: true, {12, 4}: true, {11, 4}: true, {10, 4}: true,
		{1, 3}: true, {13, 3}: true, {12, 3}: true, {11, 3}: true, {10, 3}: true,
		{1, 2}: true, {13, 2}: true, {12, 2}: true, {11, 2}: true, {10, 2}: true,
		{1, 1}: true, {13, 1}: true, {12, 1}: true, {11, 1}: true, {10, 1}: true,
	})
	hand := GetStraightFlush(deck)
	fmt.Printf("straight flush: %v\n", hand)
	hand = GetFlush(deck)
	fmt.Printf("flush: %v\n", hand)
	hand = GetStraight(deck)
	fmt.Printf("straight: %v\n", hand)
}

func TestGetMatches(t *testing.T) {
	deck := MakeDeck(map[Card]bool{
		{1, 4}: true, {13, 4}: true, {12, 4}: true, {11, 4}: true, {10, 4}: true,
		{1, 3}: true, {13, 3}: true, {12, 3}: true, {11, 3}: true, {10, 3}: true,
		{1, 2}: true, {13, 2}: true, {12, 2}: true, {11, 2}: true, {10, 2}: true,
		{1, 1}: true, {13, 1}: true, {12, 1}: true, {11, 1}: true, {10, 1}: true,
	})
	hand := GetFourOfAKind(deck)
	fmt.Printf("four of a kind: %v\n", hand)
	hand = GetFullHouse(deck)
	fmt.Printf("full house: %v\n", hand)
	hand = GetThreeOfAKind(deck)
	fmt.Printf("three of a kind: %v\n", hand)
	hand = GetTwoPair(deck)
	fmt.Printf("two pair: %v\n", hand)
	hand = GetPair(Hand{Card{13, 2}, Card{13, 4}, Card{11, 4}, Card{6, 4}, Card{2, 2}})
	fmt.Printf("pair: %v\n", hand)
}

func TestFiveHighest(t *testing.T) {
	deck := MakeDeck(map[Card]bool{
		{1, 4}: true, {13, 4}: true, {12, 4}: true, {11, 4}: true, {10, 4}: true,
		{1, 3}: true, {13, 3}: true, {12, 3}: true, {11, 3}: true, {10, 3}: true,
		{1, 2}: true, {13, 2}: true, {12, 2}: true, {11, 2}: true, {10, 2}: true,
		{1, 1}: true, {13, 1}: true, {12, 1}: true, {11, 1}: true, {10, 1}: true,
	})
	hand := GetFiveHighestCards(deck)
	fmt.Printf("highest cards: %v\n", hand)
}

func TestBestHand(t *testing.T) {
	deck := MakeDeck(map[Card]bool{
		{1, 4}: true, {13, 4}: true, {12, 4}: true, {11, 4}: true, {10, 4}: true,
		{1, 3}: true, {13, 3}: true, {12, 3}: true, {11, 3}: true, {10, 3}: true,
		{1, 2}: true, {13, 2}: true, {12, 2}: true, {11, 2}: true, {10, 2}: true,
		{1, 1}: true, {13, 1}: true, {12, 1}: true, {11, 1}: true, {10, 1}: true,
	})
	hand := BestHand(deck)
	fmt.Printf("best hand: %v\n", hand)
}
