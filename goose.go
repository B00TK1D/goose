package main

import (
	"fmt"
)

func LCS(inputs [][]byte) (int, []int) {
	inputsCount := len(inputs)
	matchIndices := make([]int, inputsCount)
	lengths := make([]int, inputsCount)
	minLen := len(inputs[0])
	for i, input := range inputs {
		lengths[i] = len(input)
		if lengths[i] < minLen {
			minLen = lengths[i]
		}
	}
	matchedLen := 0
	upperBound := minLen
	lowerBound := 1
	subsetLen := minLen
	startIndex1 := 0
	startIndex2 := 0
	for minLen > 0 {
		subsetIndex := inputsCount - 1
		subsetLen = (upperBound + lowerBound) / 2
		startIndex1 = 0
		for startIndex1 <= lengths[0]-subsetLen {
			startIndex2 = 0
			subsetIndex = inputsCount - 1
			for startIndex2 <= lengths[subsetIndex]-subsetLen && subsetIndex > 0 {
				equal := true
				for checkIndex := 0; checkIndex < subsetLen; checkIndex++ {
					if inputs[0][startIndex1+checkIndex] != inputs[subsetIndex][startIndex2+checkIndex] {
						equal = false
						break
					}
				}
				if equal {
					matchIndices[subsetIndex] = startIndex2
					subsetIndex--
					startIndex2 = 0
					continue
				}
				startIndex2++
			}
			if subsetIndex == 0 {
				matchIndices[0] = startIndex1
				matchedLen = subsetLen
				break
			}
			startIndex1++
		}
		if subsetIndex == 0 {
			lowerBound = subsetLen + 1
		} else {
			upperBound = subsetLen - 1
		}
		if lowerBound > upperBound {
			break
		}
	}
	return matchedLen, matchIndices
}

func RLCS(inputs [][]byte) int {
	preInputs := [][]byte{}
	postInputs := [][]byte{}
	matchLen, indices := LCS(inputs)
	if matchLen == 0 {
		return 0
	}
	for i, index := range indices {
		preInputs = append(preInputs, inputs[i][0:index])
		postInputs = append(postInputs, inputs[i][index+matchLen:len(inputs[i])])
	}
	return matchLen + RLCS(preInputs) + RLCS(postInputs)
}

func normalizedRLCS(inputs [][]byte) float64 {
	maxLen := 0
	for _, input := range inputs {
		length := len(input)
		if maxLen < length {
			maxLen = length
		}
	}
	return float64(RLCS(inputs)) / float64(maxLen)
}

func normalizedCombinationRLCS(combination [][][]byte) float64 {
	totalRLCS := 0.0
  elements := 0.0
	for _, cluster := range combination {
    clusterLen := float64(len(cluster))
		totalRLCS += normalizedRLCS(cluster) * clusterLen
    elements += clusterLen
	}
  return totalRLCS / elements
}

func addInputToCombinations(input []byte, combinations [][][][]byte) [][][][]byte {
	if len(combinations) == 0 {
		return [][][][]byte{[][][]byte{[][]byte{input}}}
	}
	newCombinations := [][][][]byte{}
	for _, combination := range combinations {
		for modifyGroupIndex := range combination {
			newCombination := [][][]byte{}
			for groupIndex, group := range combination {
				newCombination = append(newCombination, group)
				if groupIndex == modifyGroupIndex {
					newCombination[groupIndex] = append(newCombination[groupIndex], input)
				}
			}
			newCombinations = append(newCombinations, newCombination)
		}
		newCombinations = append(newCombinations, append(combination, [][]byte{input}))
	}
	return newCombinations
}

func bruteForceClustering(inputs [][]byte, tuning float64) [][][]byte {
	inputCount := len(inputs)
	combinations := [][][][]byte{}

	for i := range inputCount {
		combinations = addInputToCombinations(inputs[i], combinations)
	}
	bestCombination := [][][]byte{}
	bestScore := 0.0
	totalElements := 0.0
	for _, combination := range combinations {
		elements := float64(len(combination))
		if elements > totalElements {
			totalElements = elements
		}
	}
	for _, combination := range combinations {
		scoreRLCS := normalizedCombinationRLCS(combination)
		scoreClusters := (totalElements - float64(len(combination))) / totalElements
		score := (tuning * scoreRLCS) + ((1.0 - tuning) * scoreClusters)
		if score >= bestScore {
			bestCombination = combination
			bestScore = score
		}
	}
	return bestCombination
}

func printCombination(combination [][][]byte) {
	fmt.Println("Combination:")
	for _, group := range combination {
		var groupStrings []string
		for _, b := range group {
			groupStrings = append(groupStrings, string(b))
		}
		fmt.Println("  Group:", groupStrings)
	}
}

func printCombinations(combinations [][][][]byte) {
	for _, combination := range combinations {
		printCombination(combination)
	}
}

func main() {
	tests := [][]byte{
		[]byte("428efctesting123"),
		[]byte("444efgtesting456"),
		[]byte("424efgtesting456"),
	}

	//rlcsSum := normalizedRLCS(tests)
	//fmt.Println("Normalized RLCS:", rlcsSum)
	optimal := bruteForceClustering(tests, 0.5)
  optimalRLCS := normalizedCombinationRLCS(optimal)
	printCombination(optimal)
  fmt.Println("RLCS:", optimalRLCS)
}
