package main

import (
	"fmt"
	"strconv"
)

// Pelicula is a struct
type Pelicula struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
	Synopsis string `json:"synopsis"`
}

func main() {
	values := []float64{1.1, 2.2, 3.1}

	fmt.Print(arrayToString(values))
}

func arrayToString(array []float64) string {
	newArray := strconv.FormatFloat(array[0], 'f', 6, 64)
	for i := 1; i < len(array); i++ {
		newArray = newArray + "," + strconv.FormatFloat(array[i], 'f', 6, 64)
	}
	return newArray
}
