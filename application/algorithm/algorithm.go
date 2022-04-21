package algorithm

import (
	"sort"
)

// Comparisons
func CardLessThan(card1, card2 Card) bool {
	val1 := card1.Value
	val2 := card2.Value
	if val1 != 1 && val2 != 1 || val1 == 1 && val2 == 1 {
		return card1.Value < card2.Value || (card1.Value == card2.Value && card1.Suit < card2.Suit)
	} else if val2 == 1 {
		return true
	} else {
		return false
	}
}

func HandLessThan(hand1 Hand, hand2 Hand) bool {
	var best1 Hand5
	var best2 Hand5
	best1 = GetRoyalFlush(hand1)
	best2 = GetRoyalFlush(hand2)
	if best2[0].Value != 0 {
		if best1[0].Value != 0 {
			return RoyalFlushLessThan(best1, best2)
		} else {
			return false
		}
	} else if best1[0].Value != 0 {
		return true
	}
	best1 = GetStraightFlush(hand1)
	best2 = GetStraightFlush(hand2)
	if best2[0].Value != 0 {
		if best1[0].Value != 0 {
			return StraightFlushLessThan(best1, best2)
		} else {
			return false
		}
	} else if best1[0].Value != 0 {
		return true
	}
	best1 = GetFourOfAKind(hand1)
	best2 = GetFourOfAKind(hand2)
	if best2[0].Value != 0 {
		if best1[0].Value != 0 {
			return FourOfAKindLessThan(best1, best2)
		}
		return false
	} else if best1[0].Value != 0 {
		return true
	}
	best1 = GetFullHouse(hand1)
	best2 = GetFullHouse(hand2)
	if best2[0].Value != 0 {
		if best1[0].Value != 0 {
			return FullHouseLessThan(best1, best2)
		}
		return false
	} else if best1[0].Value != 0 {
		return true
	}
	best1 = GetFlush(hand1)
	best2 = GetFlush(hand2)
	if best2[0].Value != 0 {
		if best1[0].Value != 0 {
			return FlushLessThan(best1, best2)
		} else {
			return false
		}
	} else if best1[0].Value != 0 {
		return true
	}
	best1 = GetStraight(hand1)
	best2 = GetStraight(hand2)
	if best2[0].Value != 0 {
		if best1[0].Value != 0 {
			return StraightLessThan(best1, best2)
		} else {
			return false
		}
	} else if best1[0].Value != 0 {
		return true
	}
	best1 = GetThreeOfAKind(hand1)
	best2 = GetThreeOfAKind(hand2)
	if best2[0].Value != 0 {
		if best1[0].Value != 0 {
			return ThreeOfAKindLessThan(best1, best2)
		}
		return false
	} else if best1[0].Value != 0 {
		return true
	}
	best1 = GetTwoPair(hand1)
	best2 = GetTwoPair(hand2)
	if best2[0].Value != 0 {
		if best1[0].Value != 0 {
			return TwoPairLessThan(best1, best2)
		}
		return false
	} else if best1[0].Value != 0 {
		return true
	}
	best1 = GetPair(hand1)
	best2 = GetPair(hand2)
	if best2[0].Value != 0 {
		if best1[0].Value != 0 {
			return PairLessThan(best1, best2)
		}
		return false
	} else if best1[0].Value != 0 {
		return true
	}
	best1 = GetFiveHighestCards(hand1)
	best2 = GetFiveHighestCards(hand2)
	return HighCardLessThan(best1, best2)
}

func RoyalFlushLessThan(hand1, hand2 Hand5) bool {
	return hand1[0].Suit < hand2[0].Suit
}

func StraightFlushLessThan(hand1, hand2 Hand5) bool {
	return hand1[4].Value < hand2[4].Value
}

func FourOfAKindLessThan(hand1, hand2 Hand5) bool {
	if hand1[0].Value < hand2[0].Value {
		return true
	} else if hand1[0].Value == hand2[0].Value && hand1[4].Value < hand2[4].Value {
		return true
	}
	return false
}

func FullHouseLessThan(hand1, hand2 Hand5) bool {
	if hand1[0].Value < hand2[0].Value {
		return true
	} else if hand1[0].Value == hand2[0].Value && hand1[3].Value < hand2[3].Value {
		return true
	}
	return false
}

func FlushLessThan(hand1, hand2 Hand5) bool {
	return hand1[0].Value < hand2[0].Value
}

func StraightLessThan(hand1, hand2 Hand5) bool {
	return hand1[4].Value < hand2[4].Value
}

func ThreeOfAKindLessThan(hand1, hand2 Hand5) bool {
	if hand1[0].Value < hand2[0].Value {
		return true
	} else if hand1[0].Value == hand2[0].Value && hand1[3].Value < hand2[3].Value {
		return true
	}
	return false
}

func TwoPairLessThan(hand1, hand2 Hand5) bool {
	if hand1[0].Value < hand2[0].Value {
		return true
	} else if hand1[0].Value == hand2[0].Value {
		if hand1[2].Value < hand2[2].Value {
			return true
		} else if hand1[2].Value == hand2[2].Value && hand1[4].Value < hand2[4].Value {
			return true
		}
	}
	return false
}

func PairLessThan(hand1, hand2 Hand5) bool {
	if hand1[0].Value < hand2[0].Value {
		return true
	} else if hand1[0].Value == hand2[0].Value {
		if hand1[2].Value < hand2[2].Value {
			return true
		} else if hand1[2].Value == hand2[2].Value {
			if hand1[3].Value < hand2[3].Value {
				return true
			} else if hand1[3].Value == hand2[3].Value && hand1[4].Value < hand2[4].Value {
				return true
			}
		}
	}
	return false
}

func HighCardLessThan(hand1, hand2 Hand5) bool {
	for i := 0; i < 5; i++ {
		if hand1[i].Value != hand2[i].Value {
			return CardLessThan(hand1[i], hand2[i])
		}
	}
	return false
}

// Get best 5 card hand functions
func BestHand(hand Hand) Hand5 {
	var best Hand5
	if best = GetRoyalFlush(hand); best[0].Value != 0 {
		return best
	} else if best = GetStraightFlush(hand); best[0].Value != 0 {
		return best
	} else if best = GetFourOfAKind(hand); best[0].Value != 0 {
		return best
	} else if best = GetFullHouse(hand); best[0].Value != 0 {
		return best
	} else if best = GetFlush(hand); best[0].Value != 0 {
		return best
	} else if best = GetStraight(hand); best[0].Value != 0 {
		return best
	} else if best = GetThreeOfAKind(hand); best[0].Value != 0 {
		return best
	} else if best = GetTwoPair(hand); best[0].Value != 0 {
		return best
	} else if best = GetPair(hand); best[0].Value != 0 {
		return best
	} else {
		return GetFiveHighestCards(hand)
	}
}

func GetRoyalFlush(hand Hand) Hand5 {
	dict := map[Card]bool{}
	for _, card := range hand {
		dict[card] = true
	}

	if dict[Card{1, 4}] && dict[Card{13, 4}] && dict[Card{12, 4}] && dict[Card{11, 4}] && dict[Card{10, 4}] {
		return Hand5{Card{1, 4}, Card{13, 4}, Card{12, 4}, Card{11, 4}, Card{10, 4}}
	} else if dict[Card{1, 3}] && dict[Card{13, 3}] && dict[Card{12, 3}] && dict[Card{11, 3}] && dict[Card{10, 3}] {
		return Hand5{Card{1, 3}, Card{13, 3}, Card{12, 3}, Card{11, 3}, Card{10, 3}}
	} else if dict[Card{1, 2}] && dict[Card{13, 2}] && dict[Card{12, 2}] && dict[Card{11, 2}] && dict[Card{10, 2}] {
		return Hand5{Card{1, 2}, Card{13, 2}, Card{12, 2}, Card{11, 2}, Card{10, 2}}
	} else if dict[Card{1, 1}] && dict[Card{13, 1}] && dict[Card{12, 1}] && dict[Card{11, 1}] && dict[Card{10, 1}] {
		return Hand5{Card{1, 1}, Card{13, 1}, Card{12, 1}, Card{11, 1}, Card{10, 1}}
	} else {
		return Hand5{}
	}
}

func GetStraightFlush(hand Hand) Hand5 {
	valueMap := [13]Hand{}
	for _, card := range hand {
		valueMap[card.Value-1] = append(valueMap[card.Value-1], card)
	}
	suitMap := map[int]Hand{}
	for i := 13; i > 0; i-- {
		for _, card := range valueMap[i%13] {
			suitMap[card.Suit] = append(suitMap[card.Suit], card)
		}
	}

	for i := 4; i > 0; i-- {
		if len(suitMap[i]) > 4 {
			consecutive := 1
			for j := 1; j < len(suitMap[i]); j++ {
				if suitMap[i][j-1].Value == 1 && suitMap[i][j].Value == 13 {
					consecutive++
				} else if suitMap[i][j].Value == suitMap[i][j-1].Value-1 {
					consecutive++
				} else {
					consecutive = 1
					break
				}
				if consecutive == 5 {
					return Hand5{
						suitMap[i][j],
						suitMap[i][j-1],
						suitMap[i][j-2],
						suitMap[i][j-3],
						suitMap[i][j-4],
					}
				}
			}
		}
	}

	return Hand5{}
}

func GetFourOfAKind(hand Hand) Hand5 {
	matches := GetMatchesInHand(hand, map[int]int{4: 1})
	if len(matches[4]) > 0 {
		return Hand5{
			matches[4][0][0],
			matches[4][0][1],
			matches[4][0][2],
			matches[4][0][3],
			matches[1][0][0],
		}
	} else {
		return Hand5{}
	}
}

func GetFullHouse(hand Hand) Hand5 {
	matches := GetMatchesInHand(hand, map[int]int{3: 1, 2: 1})

	if len(matches[3]) > 0 && len(matches[2]) > 0 {
		return Hand5{
			matches[3][0][0],
			matches[3][0][1],
			matches[3][0][2],
			matches[2][0][0],
			matches[2][0][1],
		}
	} else {
		return Hand5{}
	}
}

func GetFlush(hand Hand) Hand5 {
	valueMap := [13]Hand{}
	for _, card := range hand {
		valueMap[card.Value-1] = append(valueMap[card.Value-1], card)
	}
	suitMap := map[int]Hand{}
	for i := 13; i > 0; i-- {
		for _, card := range valueMap[i%13] {
			suitMap[card.Suit] = append(suitMap[card.Suit], card)
		}
	}
	maxVal, maxSuit := 0, 0
	for i := 4; i > 0; i-- {
		if len(suitMap[i]) > 4 {
			if suitMap[i][0].Value == 1 {
				maxSuit = i
				break
			} else if suitMap[i][0].Value > maxVal {
				maxVal = suitMap[i][0].Value
				maxSuit = i
			}
		}
	}

	if maxSuit > 0 {
		return Hand5{
			suitMap[maxSuit][0],
			suitMap[maxSuit][1],
			suitMap[maxSuit][2],
			suitMap[maxSuit][3],
			suitMap[maxSuit][4],
		}
	} else {
		return Hand5{}
	}
}

func GetStraight(hand Hand) Hand5 {
	var valueTable [13]Hand
	for _, card := range hand {
		valueTable[card.Value-1] = append(valueTable[card.Value-1], card)
	}
	consecutive := 0
	for i := 13; i >= 0; i-- {
		if len(valueTable[i%13]) > 0 {
			consecutive++
			if consecutive == 5 {
				return Hand5{
					valueTable[i%13][0],
					valueTable[(i+1)%13][0],
					valueTable[(i+2)%13][0],
					valueTable[(i+3)%13][0],
					valueTable[(i+4)%13][0],
				}
			}
		} else {
			consecutive = 0
		}
	}
	return Hand5{}
}

func GetThreeOfAKind(hand Hand) Hand5 {
	matches := GetMatchesInHand(hand, map[int]int{3: 1})
	if len(matches[3]) > 0 {
		return Hand5{
			matches[3][0][0],
			matches[3][0][1],
			matches[3][0][2],
			matches[1][0][0],
			matches[1][1][0],
		}
	} else {
		return Hand5{}
	}
}

func GetTwoPair(hand Hand) Hand5 {
	matches := GetMatchesInHand(hand, map[int]int{2: 2})
	if len(matches[2]) > 1 {
		return Hand5{
			matches[2][0][0],
			matches[2][0][1],
			matches[2][1][0],
			matches[2][1][1],
			matches[1][0][0],
		}
	} else {
		return Hand5{}
	}
}

func GetPair(hand Hand) Hand5 {
	matches := GetMatchesInHand(hand, map[int]int{2: 1})

	if len(matches[2]) > 0 {
		return Hand5{
			matches[2][0][0],
			matches[2][0][1],
			matches[1][0][0],
			matches[1][1][0],
			matches[1][2][0],
		}
	} else {
		return Hand5{}
	}
}

// // Orders the matches by number (ie 2, 3, 4) and then value of match
func GetMatchesInHand(hand Hand, matchSizes map[int]int) map[int][]Hand {
	suitMap := [14]Hand{}
	for _, card := range hand {
		suitMap[card.Value] = append(suitMap[card.Value], card)
	}
	countMap := map[int][]Hand{}
	for i := 13; i > 0; i-- {
		count := len(suitMap[i])
		found := false
		if count == 0 {
			continue
		}
		for matchSize, matchLength := range matchSizes {
			if count >= matchSize && len(countMap[matchSize]) < matchLength {
				addHand := Hand{}
				for j, card := range suitMap[i] {
					if j < matchSize {
						addHand = append(addHand, card)
					} else {
						countMap[1] = append(countMap[1], Hand{card})
					}
				}
				countMap[matchSize] = append(countMap[matchSize], addHand)
				found = true
				break
			}
		}
		if !found {
			for _, card := range suitMap[i] {
				countMap[1] = append(countMap[1], Hand{card})
			}
		}
	}
	return countMap
}

func GetFiveHighestCards(hand Hand) Hand5 {
	sort.SliceStable(hand, func(i, j int) bool {
		return CardLessThan(hand[i], hand[j])
	})
	return Hand5{hand[4], hand[3], hand[2], hand[1], hand[0]}
}

func GetHighCard(hand Hand) Card {
	maxCard := Card{0, 0}
	for _, card := range hand {
		if CardLessThan(maxCard, card) {
			maxCard = card
		}
	}
	return maxCard
}

// general validators
func IsRoyalFlush(hand Hand) bool {
	dict := map[Card]bool{}
	for _, card := range hand {
		dict[card] = true
	}

	return (dict[Card{1, 1}] && dict[Card{13, 1}] && dict[Card{12, 1}] && dict[Card{11, 1}] && dict[Card{10, 1}]) ||
		(dict[Card{1, 2}] && dict[Card{13, 2}] && dict[Card{12, 2}] && dict[Card{11, 2}] && dict[Card{10, 2}]) ||
		(dict[Card{1, 3}] && dict[Card{13, 3}] && dict[Card{12, 3}] && dict[Card{11, 3}] && dict[Card{10, 3}]) ||
		(dict[Card{1, 4}] && dict[Card{13, 4}] && dict[Card{12, 4}] && dict[Card{11, 4}] && dict[Card{10, 4}])
}

func IsStraightFlush(hand Hand) bool {
	var valueTable [13]int
	for _, card := range hand {
		valueTable[card.Value-1] = card.Suit
	}
	consecutive := 0
	lastSuit := valueTable[0]
	for i := 0; i < 14; i++ {
		if valueTable[i%13] > 0 && valueTable[i%13] == lastSuit {
			consecutive++
			if consecutive == 5 {
				return true
			}
		} else {
			lastSuit = valueTable[i%13]
			consecutive = 0
		}
	}
	return false
}

func IsFourOfAKind(hand Hand) bool {
	return matchCountsInHand(hand)[4] > 0
}

func IsFullHouse(hand Hand) bool {
	countMap := matchCountsInHand(hand)
	return countMap[3] > 0 && countMap[2] > 0
}

func IsFlush(hand Hand) (bool, []int) {
	suitMap := map[int]int{1: 0, 2: 0, 3: 0, 4: 0}
	flushes := []int{}
	for _, card := range hand {
		suitMap[card.Suit]++
		if suitMap[card.Suit] == 5 {
			flushes = append(flushes, card.Suit)
		}
	}
	return len(flushes) > 0, flushes
}

func IsStraight(hand Hand) bool {
	var valueTable [13]bool
	for _, card := range hand {
		valueTable[card.Value-1] = true
	}
	consecutive := 0
	for i := 0; i < 14; i++ {
		if valueTable[i%13] {
			consecutive++
			if consecutive == 5 {
				return true
			}
		} else {
			consecutive = 0
		}
	}
	return false
}

func IsThreeOfAKind(hand Hand) bool {
	return matchCountsInHand(hand)[3] > 0
}

func IsTwoPair(hand Hand) bool {
	return matchCountsInHand(hand)[2] > 1
}

func IsPair(hand Hand) bool {
	return matchCountsInHand(hand)[2] > 0
}

func matchCountsInHand(hand Hand) map[int]int {
	counts := map[int]int{}
	for _, card := range hand {
		if _, inMap := counts[card.Value]; !inMap {
			counts[card.Value] = 1
		} else {
			counts[card.Value]++
		}
	}
	countMap := map[int]int{}
	for _, count := range counts {
		if _, inMap := countMap[count]; !inMap {
			countMap[count] = 1
		} else {
			countMap[count]++
		}
	}
	return countMap
}

// five card validators

func IsRoyalFlush5(hand Hand5) bool {
	isFlush, suit := IsFlush5(hand)
	if !isFlush {
		return false
	}
	dict := map[Card]bool{}
	for _, card := range hand {
		dict[card] = true
	}
	return dict[Card{1, suit}] && dict[Card{13, suit}] && dict[Card{12, suit}] && dict[Card{11, suit}] && dict[Card{10, suit}]
}

func IsStraightFlush5(hand Hand5) bool {
	if isFlush, _ := IsFlush5(hand); !isFlush {
		return false
	}
	if isStraight := IsStraight5(hand); !isStraight {
		return false
	}
	return true
}

func IsFourOfAKind5(hand Hand5) bool {
	return matchCountsInHand5(hand)[4] > 0
}

func IsFullHouse5(hand Hand5) bool {
	countMap := matchCountsInHand5(hand)
	return countMap[3] > 0 && countMap[2] > 0
}

func IsFlush5(hand Hand5) (bool, int) {
	suit := 0
	for _, card := range hand {
		if suit == 0 {
			suit = card.Suit
		} else {
			if card.Suit != suit {
				return false, suit
			}
		}
	}
	return true, suit
}

func IsStraight5(hand Hand5) bool {
	var valueTable [13]bool
	for _, card := range hand {
		valueTable[card.Value-1] = true
	}
	consecutive := 0
	for i := 0; i < 14; i++ {
		if valueTable[i%13] {
			consecutive++
			if consecutive == 5 {
				return true
			}
		} else {
			return false
		}
	}
	return false
}

func IsThreeOfAKind5(hand Hand5) bool {
	return matchCountsInHand5(hand)[3] > 0
}

func IsTwoPair5(hand Hand5) bool {
	return matchCountsInHand5(hand)[2] > 1
}

func IsPair5(hand Hand5) bool {
	return matchCountsInHand5(hand)[2] > 0
}

func matchCountsInHand5(hand Hand5) map[int]int {
	counts := map[int]int{}
	for _, card := range hand {
		if _, inMap := counts[card.Value]; !inMap {
			counts[card.Value] = 1
		} else {
			counts[card.Value]++
		}
	}
	countMap := map[int]int{}
	for _, count := range counts {
		if _, inMap := countMap[count]; !inMap {
			countMap[count] = 1
		} else {
			countMap[count]++
		}
	}
	return countMap
}
