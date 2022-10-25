package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

type Rank int

const (
	Follower Rank = iota
	Candidate
	Leader
)

func (r Rank) String() string {
	switch r {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	}
	return "unknown"
}

const (
	minPort  = 8091
	numNodes = 3
	timeout  = 1500

	url            = "http://localhost"
	api            = "/api/v1"
	endAppendLogs  = "/append-logs"
	endRequestVote = "/request-vote"
)

var (
	ping chan int
	rank Rank
)

func AppendLogs(w http.ResponseWriter, r *http.Request) {
	if rank == Follower {
		ping <- 1
	}

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
	fmt.Printf("%s:%s\n", port, string(b))
}

func HeartBeat(ownPort int) {
	for i := 0; i < numNodes; i++ {
		//Ping all ports except self
		if minPort+i != ownPort {
			SendEmptyAppendLogs(fmt.Sprintf(":%v", (minPort + i)))
		}
	}
}

// Wait for ping, if no ping received within timeout, promote self to candidate
func followerAction(ping <-chan int, rGen *rand.Rand) {
	select {
	case <-time.After(time.Duration(timeout+rGen.Intn(timeout)) * time.Millisecond): //Timeout
		fmt.Println("Promoted to Candidate")
		rank++
	case <-ping: //Pinged
		fmt.Println("Ping recieved")
	}
}

// Need to implement, automatically upgrade to leader for now
func candidateAction() {
	rank++
}

// Ping all other nodes with empty appendLogs call periodically to prevent timeouts
func leaderAction(ownPort int) {
	select {
	case <-time.After(time.Duration(timeout) * time.Millisecond):
		HeartBeat(ownPort)
	}
}

func performRankAction(ping <-chan int, ownPort int, rGen *rand.Rand) {
	fmt.Println(rank.String())
	switch rank {
	case Follower:
		followerAction(ping, rGen)
	case Candidate:
		candidateAction()
	case Leader:
		leaderAction(ownPort)
	}
}

func main() {
	var ownPort int
	ping = make(chan int, 0)
	defer close(ping)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		FindAndServePort(&ownPort)
	}()

	go func() {
		rGen := rand.New(rand.NewSource(time.Now().UnixNano()))
		for {
			performRankAction(ping, ownPort, rGen)
		}

	}()

	wg.Wait()
}
