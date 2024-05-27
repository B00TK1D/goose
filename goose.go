package main

import (
	"math"
	"math/rand"
)

func combinationScore(combination [][][]byte, objective func([][]byte) float64) float64 {
	totalScore := 0.0
	for _, group := range combination {
		totalScore += objective(group)
	}
	return totalScore / float64(len(combination))
}

func addInputToCombinations(input []byte, combinations [][][][]byte, maxNumCombinations int) [][][][]byte {
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
		if len(combination) < maxNumCombinations {
			newCombinations = append(newCombinations, append(combination, [][]byte{input}))
		}
	}
	return newCombinations
}

func BruteForceClustering(inputs [][]byte, objective func([][]byte) float64) [][][]byte {
	inputCount := len(inputs)
	combinations := [][][][]byte{}

	for i := range inputCount {
		combinations = addInputToCombinations(inputs[i], combinations, inputCount)
	}
	bestCombination := [][][]byte{}
	bestScore := 0.0
	for _, combination := range combinations {
		score := combinationScore(combination, objective)
		if score >= bestScore {
			bestCombination = combination
			bestScore = score
		}
	}
	return bestCombination
}

func HillClimbClustering(inputs [][]byte, objective func([][]byte) float64) [][][]byte {
	inputCount := len(inputs)
	bestCombination := [][][]byte{}

	bestScore := 0.0
	for combinationSize := range inputCount - 1 {
		combinations := [][][][]byte{}
		for i := range inputCount {
			combinations = addInputToCombinations(inputs[i], combinations, combinationSize+1)
		}
		bestSizeScore := 0.0
		bestSizeCombination := bestCombination
		for _, combination := range combinations {
			if len(combination) != combinationSize+1 {
				continue
			}
			score := combinationScore(combination, objective)
			if score >= bestSizeScore {
				bestSizeScore = score
				bestSizeCombination = combination
				bestCombination = combination
			}
		}
		if bestSizeScore >= bestScore {
			bestScore = bestSizeScore
			bestCombination = bestSizeCombination
		} else {
			return bestCombination
		}
	}
	return bestCombination
}

func generateRandomPopulation(parents [][][]byte, inputs [][]byte) [][][]byte {
	inputCount := len(inputs)
	population := [][][]byte{}
	for _, group := range parents {
		population = append(population, group)
	}
	// Mutate the parents
	for _, input := range inputs {
		randIndex := rand.Intn(inputCount)
		if randIndex >= len(population) {
			population = append(population, [][]byte{input})
		} else {
			population[randIndex] = append(population[randIndex], input)
		}
	}
	return population
}

func geneticSelection(population [][][]byte, objective func([][]byte) float64) ([][][]byte, [][]byte) {
	// Evaluate set of candidates
	nextPopulation := [][][]byte{}
	deadInputs := [][]byte{}
	for _, group := range population {
		groupScore := objective(group)

		// Mutation for next generation
		threshold := rand.Float64() * rand.Float64()
		if groupScore > threshold {
			nextPopulation = append(nextPopulation, group)
		} else {
			for _, value := range group {
				deadInputs = append(deadInputs, value)
			}
		}
	}
	return nextPopulation, deadInputs
}

func geneticRecombination(population [][][]byte, objective func([][]byte) float64) [][][]byte {
	bestScore := combinationScore(population, objective)
	bestPopulation := population

	for group1Index := range population {
		for group2Index, group2 := range population {
			if group1Index >= group2Index {
				continue
			}
			tmpPopulation := [][][]byte{}
			for i, group := range population {
				if i == group2Index {
					continue
				}
				tmpPopulation = append(tmpPopulation, group)
			}
			tmpPopulation[group1Index] = append(tmpPopulation[group1Index], group2...)
			tmpScore := combinationScore(tmpPopulation, objective)
			if tmpScore > bestScore {
				bestScore = tmpScore
				bestPopulation = tmpPopulation
			}
		}
	}

	if bestScore > combinationScore(population, objective) {
		return geneticRecombination(bestPopulation, objective)
	}
	return bestPopulation
}

func GeneticClustering(inputs [][]byte, objective func([][]byte) float64) [][][]byte {

	overallBestScore := 0.0
	overallBestPopulation := [][][]byte{}

	attemptCount := int(math.Log(float64(len(inputs))) + 2.0)

	for range attemptCount {
		deadInputs := [][]byte{}
		for _, input := range inputs {
			deadInputs = append(deadInputs, input)
		}

		population := [][][]byte{}
		bestPopulation := [][][]byte{}
		currentScore := 0.0
		bestScore := 0.0

		epochCount := len(inputs) * len(inputs)

		// Select initial candidate population
		population, _ = geneticSelection(bestPopulation, objective)

		for range epochCount {
			// Mutation
			population = generateRandomPopulation(population, deadInputs)

			// Fitness assessment (Evaluation)
			currentScore = combinationScore(population, objective)

			// Selection
			if currentScore > bestScore {
				bestPopulation = population
				bestScore = currentScore
			}

			// Recombination
			population, deadInputs = geneticSelection(bestPopulation, objective)
		}

		bestPopulation = geneticRecombination(bestPopulation, objective)
		bestScore = combinationScore(bestPopulation, objective)

		if bestScore > overallBestScore {
			overallBestScore = bestScore
			overallBestPopulation = bestPopulation
		}
	}

	return overallBestPopulation
}
