package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Perceptron struct {
	eta     float64
	n_inter int
	w       []float64
	errors  []int
}

func (p *Perceptron) Fit(X [][]float64, y []int) {
	auxW := make([]float64, len(X[0])+1)
	p.w = append(p.w, auxW...)
	for i := 0; i < p.n_inter; i++ {
		errors := 0
		for j := 0; j < len(X); j++ {
			update := p.eta * float64(y[j]-p.Predict(X[j]))
			fmt.Println(p.eta, y[j], p.Predict(X[j]), update)
			p.w[0] += update
			for k := 1; k < len(p.w); k++ {
				p.w[k] += update * X[j][k-1]
			}
			if update != 0.0 {
				errors += 1
			}
		}
		p.errors = append(p.errors, errors)
	}
	fmt.Println(p.w)
}

func (p *Perceptron) Predict(X []float64) int {
	if p.NetInput(X) >= 0.0 {
		//fmt.Println(1)
		return 1
	}
	//fmt.Println(-1)
	return -1
}

func (p *Perceptron) NetInput(X []float64) float64 {
	z := 0.0
	for i := 0; i < len(X); i++ {
		z += X[i] * p.w[i+1]
	}
	z += p.w[0]
	return z
}

func (p *Perceptron) Accuracy(xTest [][]float64, yTest []int) float64 {
	correctPredict := 0.0
	for i := 0; i < len(xTest); i++ {
		fmt.Println("Predict:", p.Predict(xTest[i]), "\tTrue: ", yTest[i])
		if p.Predict(xTest[i]) == yTest[i] {
			correctPredict++

		}
	}
	return correctPredict / float64(len(xTest))
}

func main() {
	r := redCSV("iris.data")

	X, y := fixData(r)

	y = whatYouWantToPredict(y, 0)

	xTrain, xTest, yTrain, yTest := splitData(X, y, 100)

	neuron := Perceptron{eta: 0.1, n_inter: 10}
	neuron.Fit(xTrain, yTrain)
	fmt.Println("Accuracy: ", neuron.Accuracy(xTest, yTest))
}

func redCSV(name string) [][]string {
	csvfile, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}

	r, err := csv.NewReader(csvfile).ReadAll()

	return r
}

func whatYouWantToPredict(targets []int, wanted int) []int {
	newTargets := make([]int, len(targets))
	for i := 0; i < len(targets); i++ {
		if targets[i] == wanted {
			newTargets[i] = -1
		} else {
			newTargets[i] = 1
		}
	}
	return newTargets

}

func splitData(X [][]float64, y []int, nTraining int) ([][]float64, [][]float64, []int, []int) {
	xTrain := make([][]float64, nTraining)
	xTest := make([][]float64, len(X)-nTraining)
	yTrain := make([]int, nTraining)
	yTest := make([]int, len(X)-nTraining)
	for i := 0; i < nTraining; i++ {
		xTrain[i] = X[i]
		yTrain[i] = y[i]
	}
	for i := nTraining; i < len(X); i++ {
		xTest[i-nTraining] = X[i]
		yTest[i-nTraining] = y[i]
	}
	return xTrain, xTest, yTrain, yTest
}

func fixData(r [][]string) ([][]float64, []int) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })
	targets := []int{}
	target := len(r[0]) - 1
	training := make([][]float64, len(r))

	keys := make(map[string]bool)
	list := []string{}

	for i, entry := range r {
		if _, value := keys[entry[target]]; !value {
			keys[entry[target]] = true
			list = append(list, entry[target])
			for number, irisType := range list {
				if entry[target] == irisType {
					targets = append(targets, number)
				}
			}

		} else {
			for number, irisType := range list {
				if entry[target] == irisType {
					targets = append(targets, number)
				}
			}
		}
		trainingRow := make([]float64, target)
		for j := 0; j < target; j++ {
			auxNum, _ := strconv.ParseFloat(r[i][j], 16)
			trainingRow[j] = auxNum
		}
		training[i] = trainingRow
	}
	return training, targets
}
