package algorithm

import (
	"gonum.org/v1/gonum/stat/combin"
)

func incrementIndices(indices []int, n int) []int {
	k := len(indices)
	i := k - 1
	newIndices := make([]int, k)
	copy(newIndices, indices)

	for i >= 0 {
		newIndices[i] = (newIndices[i] + 1) % (n - (k - i - 1))
		if newIndices[i] != 0 || i == 0 {
			break
		} else {
			i -= 1
		}
	}

	i += 1
	for i < k {
		newIndices[i] = newIndices[i-1] + 1
		i += 1
	}

	return newIndices
}

func indexCombinations(n int, k int) chan []int {
	ch := make(chan []int, combin.Binomial(n, k))
	done := k == 0

	var indices []int
	for i := 0; i < k; i++ {
		indices = append(indices, i)
	}

	var finalIndices []int
	for i := n - k; i < n; i++ {
		finalIndices = append(finalIndices, i)
	}

	go func() {
		for !done {
			comb := make([]int, k)
			copy(comb, indices)
			ch <- comb

			for i := 0; i < k; i++ {
				if indices[i] == finalIndices[i] {
					if i == k-1 {
						done = true
					}
				} else {
					break
				}
			}

			indices = incrementIndices(indices, n)
		}
		close(ch)
	}()

	return ch
}

func Combinations(array []interface{}, k int) chan []interface{} {
	n := len(array)
	total := combin.Binomial(n, k)
	ch := make(chan []interface{}, total)
	// count := 0

	go func() {
		for indexComb := range indexCombinations(n, k) {
			// go func(indexComb []int) {
			comb := make([]interface{}, k)
			for i := 0; i < k; i++ {
				comb[i] = array[indexComb[i]]
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
