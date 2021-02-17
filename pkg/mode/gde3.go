package mo

import (
	"encoding/csv"
	"math"
	"math/rand"
	"os"
	"sync"
)

// tokens is a counting semaphore use to
// enforce  a limit of 10 concurrent requests
var tokens = make(chan struct{}, 30)

// GD3 -> runs a simple multiObjective DE in the ZDT1 case
func GD3(
	wg *sync.WaitGroup,
	rankedCh chan<- Elements,
	maximumObjs chan<- []float64,
	p Params,
	evaluate func(e *Elem, M int) error,
	variant VariantFn,
	population Elements,
	f *os.File,
) {
	defer wg.Done()

	// adding to concurrent queue
	tokens <- struct{}{}
	defer f.Close()

	// var writer *csv.Writer
	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	// maximum objs found
	maxObjs := make([]float64, p.M)

	// calculates the objs of the inital population
	for i := range population {
		err := evaluate(&population[i], p.M)
		checkError(err)

		for j, obj := range population[i].Objs {
			maxObjs[j] = math.Max(maxObjs[j], obj)
		}
	}

	writeHeader(population, writer)
	writeGeneration(population, writer)

	// stores the rank[0] of each generation
	bestElems := make(Elements, 0)

	var genRankZero Elements
	var bestInGen Elements

	for g := 0; g < p.GEN; g++ {
		genRankZero, _ = FilterDominated(population)

		for i := 0; i < len(population); i++ {
			vr, err := variant.fn(
				population,
				genRankZero,
				varParams{
					currPos: i,
					DIM:     p.DIM,
					F:       p.F,
					P:       p.P,
				})
			checkError(err)

			// trial element
			trial := population[i].Copy()

			// CROSS OVER
			currInd := rand.Int() % p.DIM
			luckyIndex := rand.Int() % p.DIM

			for j := 0; j < p.DIM; j++ {
				changeProb := rand.Float64()
				if changeProb < p.CR || currInd == luckyIndex {
					trial.X[currInd] = vr.X[currInd]
				}

				if trial.X[currInd] < p.FLOOR {
					trial.X[currInd] = p.FLOOR
				}
				if trial.X[currInd] > p.CEIL {
					trial.X[currInd] = p.CEIL
				}
				currInd = (currInd + 1) % p.DIM
			}

			evalErr := evaluate(&trial, p.M)
			checkError(evalErr)

			// SELECTION
			comp := DominanceTest(&population[i].Objs, &trial.Objs)
			if comp == 1 {
				population[i] = trial.Copy()
			} else if comp == 0 {
				population = append(population, trial.Copy())
			}
		}

		population, bestInGen = ReduceByCrowdDistance(&population, p.NP)
		bestElems = append(bestElems, bestInGen...)

		writeGeneration(population, writer)

		// checks for the biggest objective
		for _, p := range population {
			for j, obj := range p.Objs {
				maxObjs[j] = math.Max(maxObjs[j], obj)
			}
		}
	}

	// sending via channel the data
	rankedCh <- bestElems
	maximumObjs <- maxObjs

	// clearing concurrent queue
	<-tokens
}