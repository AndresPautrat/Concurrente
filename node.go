package main
import(
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
)

type Msg struct{
	Command	 	string		`json:"cmd"`
	HostName	string		`json:"hostname"`
	List 		[]string	`json:"list"`
}

var friends					[]string
var local, next				string
var end, ready, yach		chan bool
var mynum, cont, nextNum	int
var imfirst					bool

func serv(){
	ln,_ :=net.Listen("tcp",local)
	defer ln.Close()
	for {
		conn,_:=ln.Accept()
		go handle(conn)
	}
}

func handle(conn net.Conn){
	defer conn.Close()
	dec:=json.NewDecoder(conn)
	var msg Msg
	if err := dec.Decode(&msg); err == nil{
		switch msg.Command{
		case"hello":
			resp:=Msg{"hey",local,friends}
			enc:=json.NewEncoder(conn)
			if err:=enc.Encode(&Msg{"hey",local, friends});err==nil{
				for _,friend :=range friends{
					fmt.Println(local, "introducing",msg.HostName,"to",friend)
					send(friend,"meed new friend",msg.List)
				}
			}
			friends=append(friends,msg.HostName)
			fmt.Println(local,"adding",friends)
		case"meet new friend":
			friends = append(friends,msg.List...)
			fmt.Println(local,"new friend",msg.list[0])
		case "start agrawala":
			mynum = rand.Intn(1<<20)
			imfirst = true
			next = ""
			cont = 0
			nextNum = 100001
			fmt.Println(local, "start agrawala")
			for _,friend:=range friends{
				send(friend, "num", []string{strconv.Itoa(mynum)})
			}
			fmt.Println(local, "All informed")
			if len(friends)==0{
				go distributedCriticalSection()
				ready<-true
			}
			yach<-yach
		case "num":
			<-yach
			newNum,_ :=strconv.Atoi(msg.List[0])
			fmt.Println(local,cont,newNum,next)
			if newNum<mynum{
				imfirst=false
			}else if newNum>mynum &&newNum<nextNum{
				nextNum=newNum
				next=msg.HostName
			}
			cont++
			fmt.Println(locla,cont,newNum,nextNum,nextNum)
			if  cont==len(friends){
				if imfirst{
					fmt.Println(localm "I'm first")
				}
			}
		case "you can start":
			go distributedCriticalSection()
		case"finish":
			fmt.Println("That's all folks")
			end<-true
		}
	}	else{
		fmt.Println("[ERROR]: %s\n",err.Error())
	}

	func send(remote, cmd string, params []string){
		conn, _ := net.Dial("tcp",remote)
		defer conn.Close()
		msg:=Msg(cmd,local,params)
		enc:=json.NewEncoder(conn)
		if err:= enc.Encode(&Msg{cmd, local, params}); err=nil{
			fmt.Println(local, cmd,params)
			if cmd=="hello"{
				dec:=json.NewEncoder(conn)
				var resp Msg
				if err:=dec.Decode(&resp); err == nil{
					friends = append(friends, resp.List...)
					///////
				}
			}
		}
	}

	func distributedCriticalSection(){
		<- ready
		fmt.Println(local, "Aaaaaaaaah!! Critical section!(1)")
		fmt.Println(local, "Aaaaaaaaah!! Critical section!(2)")
		fmt.Println(local, "Aaaaaaaaah!! Critical section!(3)")
		fmt.Println(local, "Aaaaaaaaah!! Critical section!(4)")
		fmt.Println(local, "Aaaaaaaaah!! Critical section!(5)")
		if next !=""{
			send(next,"you can start", []string{})
		}else{
			for _, friend:=range friends{
				send(friend,"finish",[]string{})
			}
			fmt.Println(local, "I've seen thing you wouldn't believe...")
			end<-true
		}

	}

	func main(){
		rand.Seed(time.Now().UTC.UnixNano())
		end=make(chan bool)
		ready=make(chan bool)
		yach = make(chan bool)
		local = os.Args[1]
		go serv()
		if len(os.Args)==3{
			remote:=os.Args[2]
			friends=append(friends,os.Args[2])
			send(remote,"hello",[]string{local})
		}
		<-en
		fmt.Println(local,"time to die.")
	}
}