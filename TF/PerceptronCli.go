package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
)

type Data struct {
	SepalL float64 `json:"sepal_length"`
	SepalW float64 `json:"sepal_width"`
	PetalL float64 `json:"petal_length"`
	PetalW float64 `json:"petal_width"`
	Class  string  `json:"class"`
}

func readJSON() []Data {
	data, _ := ioutil.ReadFile("irisJson.json")
	var iris []Data
	_ = json.Unmarshal(data, &iris)
	return iris
	// fmt.Printf("%v %T\n", iris, iris)
}

func SplitData(allData []Data) ([][]float64, []int) {

	fmt.Println(allData[0].Class)
	x := make([][]float64, len(allData))
	y := []int{}
	keys := make(map[string]bool)
	list := []string{}

	for i, entry := range allData {
		if _, value := keys[entry.Class]; !value {
			keys[entry.Class] = true
			list = append(list, entry.Class)
			for i, irisType := range list {
				if entry.Class == irisType {
					y = append(y, i)
				}
			}
		} else {
			for i, irisType := range list {
				if entry.Class == irisType {
					y = append(y, i)
				}
			}
		}

		xRow := make([]float64, 4)
		for j := 0; j < 4; j++ {
			listOfColumns := [4]float64{allData[i].PetalL, allData[i].PetalW, allData[i].SepalL, allData[i].SepalW}
			xRow[j] = listOfColumns[j]
		}
		x[i] = xRow
	}
	return x, y
}

func server(hostname, remote string, end chan bool) {
	ln, _ := net.Listen("tcp", hostname)
	defer ln.Close()
	fmt.Println("Listening!")
	for {
		con, _ := ln.Accept()
		handle(con, hostname, remote, end)
	}
}

func handle(con net.Conn, hostname, remote string, end chan bool) {
	//allData := readJSON()
	defer con.Close()
	dec := json.NewDecoder(con)
	var msg Msg
	print(con)
	if err := dec.Decode(&msg); err == nil {
		fmt.Printf("Message: %v\n", msg)
		send(hostname, remote, "hola")
	}
}

type Msg struct {
	Addr   string `json:"addr"`
	Option string `json:"option"`
}

func send(local, remote string, msg string) {
	if remote != "0" {
		con, _ := net.Dial("tcp", remote)
		defer con.Close()
		enc := json.NewEncoder(con)
		if err := enc.Encode(Msg{local, msg}); err == nil {
			fmt.Printf("Sending %s to %s\n", msg, remote)
		}
	}
}

func main() {
	var hostname, remote string
	var test string

	//allData := readJSON()
	//x, y := SplitData(allData)

	fmt.Print("Hostname: ")
	fmt.Scanf("%s", &hostname)
	fmt.Scanf("%s", &test)
	fmt.Print(test)
	fmt.Print("Remote hostname: ")
	fmt.Scanf("%s", &remote)

	remote = fmt.Sprintf("localhost:800%s", remote)

	hostname = fmt.Sprintf("localhost:800%s", hostname)
	end := make(chan bool)

	go server(hostname, remote, end)
	<-end
	fmt.Println("That's  all, folks!")
}
