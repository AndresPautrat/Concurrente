package main

import (
	"fmt"
)

func prueba(c chan []float64) {
	fmt.Println(10)
	aux := []float64{0.0, 0.1}
	for _, i := range aux {
		fmt.Println(i)
	}
	c <- aux
}

func main() {
	auxW := make([]float64, 2)
	nThreads := 2
	chans := make([]chan []float64, nThreads)
	for i := range chans {
		chans[i] = make(chan []float64)
	}
	for i := 0; i < nThreads; i++ {
		go prueba(chans[i])
	}
	for i := 0; i < nThreads; i++ {
		subWeights := <-chans[i]
		for j, newW := range subWeights {
			auxW[j] += newW
		}
	}
}

