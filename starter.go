package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

func send(num int, remote string){
	if remote != ""{
		con,_ := net.Dial("tcp", remote)
		defer con.Close()
		fmt.Fprintf("%d\n",num)
	}
}

func main(){
	mins := []int {3,1,4,2,-1}
	for _,num :=range nums{
		send(num,"localhost:8000")
	}
}