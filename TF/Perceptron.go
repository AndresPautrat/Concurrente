//sudo lsof -i -P -n
//sudo fuser -k Port_Number/tcp
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

// Data iris
type Data struct {
	SepalL float64 `json:"sepal_length"`
	SepalW float64 `json:"sepal_width"`
	PetalL float64 `json:"petal_length"`
	PetalW float64 `json:"petal_width"`
	Class  string  `json:"class"`
}

// perceptron
type Perceptron struct {
	eta     float64   `json:"eta"`
	n_inter int       `json:"nInter"`
	w       []float64 `json:"w"`
	errors  []int     `json:"errors"`
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
func (p *Perceptron) internalPredict(X []float64, w []float64) int {
	if p.internalNetInput(X, w) >= 0.0 {
		return 1
	}
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
		return 1
	}
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

func (p *Perceptron) GetW() []float64 {
	return p.w
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

func readJSON() []Data {
	data, _ := ioutil.ReadFile("irisJson.json")
	var iris []Data
	_ = json.Unmarshal(data, &iris)
	return iris
	// fmt.Printf("%v %T\n", iris, iris)
}

func SplitData(allData []Data) ([][]float64, []int) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allData), func(i, j int) { allData[i], allData[j] = allData[j], allData[i] })
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

func server(hostname string, end []chan bool) {
	ln, _ := net.Listen("tcp", hostname)
	defer ln.Close()
	fmt.Println("Listening!")
	for {
		con, _ := ln.Accept()
		handle(con, hostname, end)
	}
}

//TODO crear en handle fits con distintos tamaños de data
func handle(con net.Conn, hostname string, end []chan bool) {
	defer con.Close()
	dec := json.NewDecoder(con)
	var msg Msg
	print(con)
	if err := dec.Decode(&msg); err == nil {
		fmt.Printf("Message: %v\n", msg)

		if msg.Option == "msg" {
			sendPerceptron(hostname, msg.Addr)
		}
		if msg.Option == "per" {
			perAux := stringToArray(msg.Message)
			fmt.Print(perAux)
			fmt.Println("MakingTrue %d", nWaits)
			end[nWaits] <- true
			nWaits++
		}
	} else {
		fmt.Println("ErrorH: ", err)
	}
	//end <- true
}

type Msg struct {
	Addr    string `json:"addr"`
	Option  string `json:"option"`
	Message string `json:"message"`
}

func send(local, remote string, msg string) {
	if remote != "0" {
		con, _ := net.Dial("tcp", remote)
		defer con.Close()
		enc := json.NewEncoder(con)
		if err := enc.Encode(Msg{local, "msg", msg}); err == nil {
			fmt.Printf("Sending %s to %s\n", msg, remote)
		} else {
			fmt.Println("ErrorS: ", err)
		}
	}
}

func sendPerceptron(local, remote string) {
	allData := readJSON()
	x, y := SplitData(allData)
	y = whatYouWantToPredict(y, 0)
	neuron := Perceptron{eta: 0.1, n_inter: 50}
	neuron.Fit(x, y, 4)

	if remote != "0" {
		con, _ := net.Dial("tcp", remote)
		defer con.Close()
		enc := json.NewEncoder(con)
		if err := enc.Encode(Msg{local, "per", arrayToString(neuron.GetW())}); err == nil {
			fmt.Printf("Sending %s to %s\n", neuron.GetW(), remote)
		} else {
			fmt.Println("ErrorSend: ", err)
		}
	}
}

func arrayToString(array []float64) string {
	newArray := strconv.FormatFloat(array[0], 'f', 6, 64)
	for i := 1; i < len(array); i++ {
		newArray = newArray + "," + strconv.FormatFloat(array[i], 'f', 6, 64)
	}
	return newArray
}

func stringToArray(text string) []float64 {
	stringsArray := strings.Split(text, ",")
	var newArray []float64
	for _, element := range stringsArray {
		floatAux, _ := strconv.ParseFloat(element, 64)
		newArray = append(newArray, floatAux)
	}
	return newArray
}

var nWaits = 0

func main() {

	var hostname string
	var remote []string
	var test string
	fmt.Print("Hostname: ")
	fmt.Scanf("%s", &hostname)
	fmt.Scanf("%s", &test)
	fmt.Print(test)

	if hostname == "0" {
		hostname = fmt.Sprintf("localhost:800%s", hostname)
		var nConnections int
		fmt.Print("Numero de distribuciones: ")
		fmt.Scanf("%d", &nConnections)
		for i := 1; i <= nConnections; i++ {
			remoteAux := fmt.Sprintf("localhost:800%s", strconv.Itoa(i))
			remote = append(remote, remoteAux)
			fmt.Print(remote)
		}
		end := make([]chan bool, nConnections)
		go server(hostname, end)

		for _, port := range remote {
			send(hostname, port, "hola")
		}

		for i := 0; i < nConnections; i++ {
			<-end[i]
		}
		fmt.Printf("termine")
	} else if hostname > "0" {
		hostname = fmt.Sprintf("localhost:800%s", hostname)

		end := make([]chan bool, 1)
		go server(hostname, end)
		<-end[0]
	}
}
