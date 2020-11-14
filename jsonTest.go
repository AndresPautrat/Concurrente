package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Pelicula struct {
	Title    string `json:"title"`
	director string `json:"director"`
	Year     int    `json:"year"`
	Synopsis string `json:"synopsis"`
}

func main() {
	fmt.Println("Json example")
	peliculas := []Pelicula{
		{"Wonder", "Patty", 2017, "BAM!"},
		{"Joker", "Todd", 2019, "hahaha"}}
	bytes, _ := json.MarshalIndent(peliculas, "", "\t")
	fmt.Println(string(bytes))

	var peliculas2 []Pelicula
	_ = json.Unmarshal(bytes, &peliculas2)
	fmt.Println(peliculas2)

	enc := json.NewEncoder(os.Stdin)
	enc.Encode(peliculas)

	dec := json.NewDecoder(os.Stdin)
	dec.Decode(&peliculas2)

	fmt.Printf("%v %T\n", peliculas2, peliculas2)
}
