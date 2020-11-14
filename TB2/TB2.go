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

func (p *Perceptron) ConcurrentForPerceptron(subX [][]float64, subY []int, c chan []float64) {
	auxW := make([]float64, len(subX[0])+1)
	for i := 0; i < p.n_inter; i++ {
		errors := 0
		for j := 0; j < len(subX); j++ {
			update := p.eta * float64(subY[j]-p.internalPredict(subX[j], auxW))
			auxW[0] += update
			for k := 1; k < len(auxW); k++ {
				auxW[k] += update * subX[j][k-1]
			}
			if update != 0.0 {
				errors += 1
			}
		}

		p.errors = append(p.errors, errors)
	}
	c <- auxW
}

//Split the data, in litle subsets
func (p *Perceptron) Fit(X [][]float64, y []int, nThreads int) {
	auxW := make([]float64, len(X[0])+1)
	p.w = append(p.w, auxW...)
	subSetLen := int(len(X) / nThreads)
	chans := make([]chan []float64, nThreads)
	for i := range chans {
		chans[i] = make(chan []float64)
	}
	for i := 0; i < nThreads; i++ {
		go p.ConcurrentForPerceptron(X[i*subSetLen:(i+1)*subSetLen], y[i*subSetLen:(i+1)*subSetLen], chans[i])
	}
	for i := 0; i < nThreads; i++ {
		subWeights := <-chans[i]
		for j, newW := range subWeights {
			p.w[j] += newW
		}
	}
}

//Dont split the data, split the epochs
/*func (p *Perceptron) Fit2(X [][]float64, y []int, nThreads int) {
	auxW := make([]float64, len(X[0])+1)
	p.w = append(p.w, auxW...)
	p.n_inter = p.n_inter / nThreads
	chans := make([]chan []float64, nThreads)
	for i := range chans {
		chans[i] = make(chan []float64)
	}
	for i := 0; i < nThreads; i++ {
		go p.ConcurrentForPerceptron(X, y, chans[i])
	}
	for i := 0; i < nThreads; i++ {
		subWeights := <-chans[i]
		for j, newW := range subWeights {
			p.w[j] += newW
		}
	}
}
*/

func (p *Perceptron) internalPredict(X []float64, w []float64) int {
	if p.internalNetInput(X, w) >= 0.0 {
		//fmt.Println(1)
		return 1
	}
	//fmt.Println(-1)
	return -1
}

func (p *Perceptron) internalNetInput(X []float64, w []float64) float64 {
	z := 0.0
	for i := 0; i < len(X); i++ {
		z += X[i] * w[i+1]
	}
	z += w[0]
	return z
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
		//fmt.Println("Predict:", p.Predict(xTest[i]), "\tTrue: ", yTest[i])
		if p.Predict(xTest[i]) == yTest[i] {
			correctPredict++

		}
	}
	return correctPredict / float64(len(xTest))
}

func main() {
	r := redCSV("./iris.data")

	X, y := fixData(r)

	y = whatYouWantToPredict(y, 0)

	xTrain, xTest, yTrain, yTest := splitData(X, y, 100)

	neuron := Perceptron{eta: 0.1, n_inter: 50}
	neuron.Fit(xTrain, yTrain, 4)
	fmt.Println("Predict:", neuron.Predict(xTest[0]), "\tTrue: ", yTest[0])
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
	return X[0:nTraining], X[nTraining:len(X)], y[0:nTraining], y[nTraining:len(X)]
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
