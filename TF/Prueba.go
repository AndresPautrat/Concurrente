package main

import (
	"fmt"
	"strings"
)

// Pelicula is a struct
type Pelicula struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
	Synopsis string `json:"synopsis"`
}

func main() {
	values := [3]int{1.2, 2.2, 3.2}
	print
	result1 := strings.Trim(strings.Replace(fmt.Sprint(values), " ", ",", -1), "[]")
	// Convert string slice to string.
	// ... Has comma in between strings.
	fmt.Println(result1[0])

	// ... Use no separator.
	result2 := strings.Split(result1, ",")
	fmt.Println(result2)
}
