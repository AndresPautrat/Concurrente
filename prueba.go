package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

var min int

func server(hostname, remote string, end chan bool){
	ln,_:=net.Listen("tcp",hostname)
	defer ln.Close()
	for{
		con,_ :=ln.Accept()
		go handle(con,remote,end)
	}
}

func send(num int, remote string){
	if remote != ""{
		con,_ := net.Dial("tcp", remote)
		defer con.Close()
		fmt.Fprintf("%d\n",num)
	}
}

func handle(con net.Coon, remote string, end chan bool){
	defer con.Close()
	r := bufio.NewReader(con)
	msg, _ := r.ReadString('\n')
	num,_:=strconv.Atoi(strings.TrimSpace(msg))
	if min == -1{
		min=num
	}
	else if num == -1{
		send(num, remote)
		fmt.Print("%s -> %d\n",con.LocalAddr(),min)
		end<-true
	} else if num <min{
		send(min,remote)
		min=num
	} else {
		send(num, remote)
	}

}

func main() {
	var hostname, remote string
	fmt.Print("Hostname: ")
	fmt.Scanf("%s",&hostname)
	fmt.Print("Remote hostname: ")
	fmt.Scanf("%s",&remote)

	end := make(chan bool)
	min = -1
}
