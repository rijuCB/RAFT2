package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

const (
	minPort  = 8091
	numNodes = 3

	url            = "http://localhost"
	api            = "/api/v1"
	endAppendLogs  = "/append-logs"
	endRequestVote = "/request-vote"
)

func AppendLogs(w http.ResponseWriter, r *http.Request) {
	//specify status code
	w.WriteHeader(http.StatusOK)

	//update response writer
	fmt.Fprintf(w, "API is up and running")
}

func RequestVote(w http.ResponseWriter, r *http.Request) {
	//specify HTTP status code
	w.WriteHeader(http.StatusOK)

	//update response
	fmt.Fprintf(w, "You've got my vote")
}

func FindAndServePort(ownPort *int) {
	// Colours!!
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	//create a new router
	router := mux.NewRouter()

	// specify endpoints, handler functions and HTTP method
	router.HandleFunc(api+endAppendLogs, AppendLogs).Methods("POST")
	router.HandleFunc(api+endRequestVote, RequestVote).Methods("GET")
	http.Handle("/", router)

	var l net.Listener
	var err error
	//Search for open ports between 8091 & 8100
	for i := 0; i < numNodes; i++ {
		attemptPort := fmt.Sprintf(":%v", (minPort + i))
		fmt.Printf("Attempting to open port %v ... ", cyan(attemptPort))
		l, err = net.Listen("tcp", attemptPort) //Attempt to open port
		if err == nil {
			*ownPort = minPort + i
			fmt.Printf(green("success!\n"))
			break
		}
		log.Printf("%v\n%v\n", red("failed!"), yellow(err))
	}

	if l == nil {
		log.Printf(red("Unable to open port"))
		return
	}

	if err := http.Serve(l, router); err != nil { //Respond to requests
		log.Printf("ERROR!\n%v\n", err)
	}
}

func SendEmptyAppendLogs(port string) {
	resp, err := http.Post(url+port+api+endAppendLogs, "application/json", strings.NewReader(""))
	if err != nil {
		log.Println(err)
		return
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	println(string(b))
}

func HeartBeat(ownPort int) {
	for i := 0; i < numNodes; i++ {
		//Ping all ports except self
		if minPort+i != ownPort {
			SendEmptyAppendLogs(fmt.Sprintf(":%v", (minPort + i)))
		}
	}
}

func main() {
	var ownPort int

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		FindAndServePort(&ownPort)
	}()

	HeartBeat(ownPort)

	wg.Wait()
}
