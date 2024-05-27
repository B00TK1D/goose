package main

import (
	"math"
	"math/rand"
)

func combinationScore[T any](combination [][]T, objective func([]T) float64) float64 {
	totalScore := 0.0
	for _, group := range combination {
		totalScore += objective(group)
	}
	return totalScore / float64(len(combination))
}

func addInputToCombinations[T any](input T, combinations [][][]T, maxNumCombinations int) [][][]T {
	if len(combinations) == 0 {
		return [][][]T{[][]T{[]T{input}}}
	}
	newCombinations := [][][]T{}
	for _, combination := range combinations {
		for modifyGroupIndex := range combination {
			newCombination := [][]T{}
			for groupIndex, group := range combination {
				newCombination = append(newCombination, group)
				if groupIndex == modifyGroupIndex {
					newCombination[groupIndex] = append(newCombination[groupIndex], input)
				}
			}
			newCombinations = append(newCombinations, newCombination)
		}
		if len(combination) < maxNumCombinations {
			newCombinations = append(newCombinations, append(combination, []T{input}))
		}
	}
	return newCombinations
}

func BruteForceClustering[T any](inputs []T, objective func([]T) float64) [][]T {
	inputCount := len(inputs)
	combinations := [][][]T{}

	for i := range inputCount {
		combinations = addInputToCombinations(inputs[i], combinations, inputCount)
	}
	bestCombination := [][]T{}
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

func HillClimbClustering[T any](inputs []T, objective func([]T) float64) [][]T {
	inputCount := len(inputs)
	bestCombination := [][]T{}

	bestScore := 0.0
	for combinationSize := range inputCount - 1 {
		combinations := [][][]T{}
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

func generateRandomPopulation[T any](parents [][]T, inputs []T) [][]T {
	inputCount := len(inputs)
	population := [][]T{}
	for _, group := range parents {
		population = append(population, group)
	}
	// Mutate the parents
	for _, input := range inputs {
		randIndex := rand.Intn(inputCount)
		if randIndex >= len(population) {
			population = append(population, []T{input})
		} else {
			population[randIndex] = append(population[randIndex], input)
		}
	}
	return population
}

func geneticSelection[T any](population [][]T, objective func([]T) float64) ([][]T, []T) {
	// Evaluate set of candidates
	nextPopulation := [][]T{}
	deadInputs := []T{}
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

func geneticRecombination[T any](population [][]T, objective func([]T) float64) [][]T {
	bestScore := combinationScore(population, objective)
	bestPopulation := population

	for group1Index := range population {
		for group2Index, group2 := range population {
			if group1Index >= group2Index {
				continue
			}
			tmpPopulation := [][]T{}
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

func GeneticClustering[T any](inputs []T, objective func([]T) float64) [][]T {

	overallBestScore := 0.0
	overallBestPopulation := [][]T{}

	attemptCount := int(math.Log(float64(len(inputs))) + 2.0)

	for range attemptCount {
		deadInputs := []T{}
		for _, input := range inputs {
			deadInputs = append(deadInputs, input)
		}

		population := [][]T{}
		bestPopulation := [][]T{}
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
