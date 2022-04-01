package simulate

import (
	"log"

	"github.com/carlschader/poker-go-api/application/algorithm"
	"gonum.org/v1/gonum/stat/combin"
)

type SimulationResult struct {
	Hand          string  `json:"hand" bson:"hand"`
	EndHandSize   int32   `json:"end_hand_size" bson:"end_hand_size"`
	ExpectedRank  float32 `json:"expected_rank" bson:"expected_rank"`
	RoyalFlush    float32 `json:"royal_flush" bson:"royal_flush"`
	StraightFlush float32 `json:"straight_flush" bson:"straight_flush"`
	FourOfAKind   float32 `json:"four_of_a_kind" bson:"four_of_a_kind"`
	FullHouse     float32 `json:"full_house" bson:"full_house"`
	Flush         float32 `json:"flush" bson:"flush"`
	Straight      float32 `json:"straight" bson:"straight"`
	ThreeOfAKind  float32 `json:"three_of_a_kind" bson:"three_of_a_kind"`
	TwoPair       float32 `json:"two_pair" bson:"two_pair"`
	Pair          float32 `json:"pair" bson:"pair"`
	HighCard      float32 `json:"high_card" bson:"high_card"`
}

type SimulationData struct {
	RankSum       int64
	RoyalFlush    int64
	StraightFlush int64
	FourOfAKind   int64
	FullHouse     int64
	Flush         int64
	Straight      int64
	ThreeOfAKind  int64
	TwoPair       int64
	Pair          int64
	HighCard      int64
}

func SimulateHand(hand algorithm.Hand, endHandSize int, exclude map[algorithm.Card]bool, ranks map[string]int64) (*SimulationResult, error) {
	exclusionSet := map[algorithm.Card]bool{}

	for _, card := range hand {
		exclusionSet[card] = true
	}
	for card := range exclude {
		exclusionSet[card] = true
	}

	deck := algorithm.MakeDeck(exclusionSet)
	data := SimulationData{}
	i := 0
	for comb := range algorithm.CardCombinations(deck, endHandSize-len(hand)) {
		currentHand := algorithm.Hand{}
		currentHand = append(currentHand, hand...)
		currentHand = append(currentHand, comb...)
		best := algorithm.Hand5{}

		if best = algorithm.GetRoyalFlush(currentHand); best[0].Value != 0 {
			data.RoyalFlush++
		} else if best = algorithm.GetStraightFlush(currentHand); best[0].Value != 0 {
			data.StraightFlush++
		} else if best = algorithm.GetFourOfAKind(currentHand); best[0].Value != 0 {
			data.FourOfAKind++
		} else if best = algorithm.GetFullHouse(currentHand); best[0].Value != 0 {
			data.FullHouse++
		} else if best = algorithm.GetFlush(currentHand); best[0].Value != 0 {
			data.Flush++
		} else if best = algorithm.GetStraight(currentHand); best[0].Value != 0 {
			data.Straight++
		} else if best = algorithm.GetThreeOfAKind(currentHand); best[0].Value != 0 {
			data.ThreeOfAKind++
		} else if best = algorithm.GetTwoPair(currentHand); best[0].Value != 0 {
			data.TwoPair++
		} else if best = algorithm.GetPair(currentHand); best[0].Value != 0 {
			data.Pair++
		} else {
			best = algorithm.GetFiveHighestCards(currentHand)
			data.HighCard++
		}

		data.RankSum += ranks[best.Hash()]
		i++
	}

	total := combin.Binomial(len(deck), endHandSize-len(hand))
	log.Printf("totals: %d %d\n", total, i)
	result := SimulationResult{
		Hand:          hand.Hash(),
		EndHandSize:   int32(endHandSize),
		ExpectedRank:  float32(data.RankSum) / float32(total),
		RoyalFlush:    float32(data.RoyalFlush) / float32(total),
		StraightFlush: float32(data.StraightFlush) / float32(total),
		FourOfAKind:   float32(data.FourOfAKind) / float32(total),
		FullHouse:     float32(data.FullHouse) / float32(total),
		Flush:         float32(data.Flush) / float32(total),
		Straight:      float32(data.Straight) / float32(total),
		ThreeOfAKind:  float32(data.ThreeOfAKind) / float32(total),
		TwoPair:       float32(data.TwoPair) / float32(total),
		Pair:          float32(data.Pair) / float32(total),
		HighCard:      float32(data.HighCard) / float32(total),
	}

	return &result, nil
}
