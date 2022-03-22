package algorithm

import (
	"fmt"
	"strings"
	"testing"
)

func TestIncrementIndices(t *testing.T) {
	indices := []int{1, 2, 3, 4, 7}
	indices = incrementIndices(indices, 10)
	fmt.Printf("%v\n", indices)
}

func TestIndexCombinations(t *testing.T) {
	n, k, i := 7, 5, 1

	for comb := range indexCombinations(n, k) {
		fmt.Printf("%d %v\n", i, comb)
		i++
	}
}

func TestCombinations(t *testing.T) {
	array := []interface{}{"r", "g", "b", "p", "y"}
	set := map[string]bool{}
	i := 1
	k := 3

	for comb := range Combinations(array, k) {
		var stringSlice []string
		for j := 0; j < k; j++ {
			stringSlice = append(stringSlice, comb[j].(string))
		}
		set[strings.Join(stringSlice, ",")] = true
		fmt.Printf("%d %v\n", i, comb)
		i++
	}

	fmt.Printf("set size: %d\n", len(set))
}
