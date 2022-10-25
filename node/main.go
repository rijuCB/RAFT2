package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

// Rank ENUM
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

	url              = "http://localhost"
	api              = "/api/v1"
	endAppendLogs    = "/append-logs"
	paramAppendLogs  = ""
	endRequestVote   = "/request-vote"
	paramRequestVote = "/{term:[0-9]+}/{requester:[0-9]+}"
)

type Node struct {
	ping chan int
	rank Rank

	term int //Add mutex
	vote int
}

var (
	// Colours!!
	cyan   = color.New(color.FgCyan).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

func (node *Node) AppendLogs(w http.ResponseWriter, r *http.Request) {
	if node.rank == Follower {
		node.ping <- 1
	}

	//specify status code
	w.WriteHeader(http.StatusOK)

	//update response writer
	fmt.Fprintf(w, "API is up and running")
}

func (node *Node) parseVoteRequest(r *http.Request) (int, int) {
	vars := mux.Vars(r)
	term, err := strconv.Atoi(vars["term"])
	if err != nil {
		fmt.Println(red(err))
		return -1, -1
	}
	requester, err := strconv.Atoi(vars["requester"])
	if err != nil {
		fmt.Println(red(err))
		return -1, -1
	}
	return term, requester
}

func (node *Node) RequestVote(w http.ResponseWriter, r *http.Request) {
	newTerm, requester := node.parseVoteRequest(r)
	if newTerm < 0 || requester < 0 {
		fmt.Fprintf(w, red("Invalid parameters!"))
		return //Invalid values
	}

	//specify HTTP status code
	w.WriteHeader(http.StatusOK)

	fmt.Printf("Received vote request from node:%v for term %v\n", cyan(requester), yellow(newTerm))

	if newTerm > node.term {
		node.vote = requester
		node.term = newTerm
	}
	if node.rank == Follower {
		node.ping <- 1
	}

	fmt.Printf("Vote cast for %v\n", green(node.vote))

	//update response
	fmt.Fprintf(w, "%v", node.vote)
}

func (node *Node) FindAndServePort(ownPort *int) {
	//create a new router
	router := mux.NewRouter()

	// specify endpoints, handler functions and HTTP method
	router.HandleFunc(api+endAppendLogs+paramAppendLogs, node.AppendLogs).Methods("POST")
	router.HandleFunc(api+endRequestVote+paramRequestVote, node.RequestVote).Methods("GET")
	http.Handle("/", router)

	var l net.Listener
	var err error
	//Search for open ports
	for i := 0; i < numNodes; i++ {
		attemptPort := fmt.Sprintf(":%v", (minPort + i))
		fmt.Printf("Attempting to open port %v ... ", cyan(attemptPort))
		l, err = net.Listen("tcp", attemptPort) //Attempt to open port
		if err == nil {
			*ownPort = minPort + i
			fmt.Printf(green("success!\n"))
			break
		}
		fmt.Printf("%v\n%v\n", red("failed!"), yellow(err))
	}

	if l == nil {
		fmt.Printf(red("Unable to open port"))
		return
	}

	if err := http.Serve(l, router); err != nil { //Respond to requests
		fmt.Printf("ERROR!\n%v\n", err)
	}
}

func (node *Node) SendEmptyAppendLogs(endpoint string) {
	resp, err := http.Post(endpoint, "application/json", strings.NewReader(""))
	if err != nil {
		fmt.Println(red(err))
		return
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(red(err))
		return
	}
	fmt.Printf("%s\n%s\n", cyan(endpoint), green(string(b)))
}

func (node *Node) HeartBeat(ownPort int) {
	for i := 0; i < numNodes; i++ {
		//Ping all ports except self
		if minPort+i != ownPort {
			node.SendEmptyAppendLogs(fmt.Sprintf("%s:%v%s%s", url, (minPort + i), api, endAppendLogs))
		}
	}
}

func (node *Node) requestVoteFromNode(endpoint string) int {
	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println(red(err))
		return -1
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(red(err))
		return -1
	}

	ballot, err := strconv.Atoi(string(b))
	if err != nil {
		fmt.Println(red(err))
		return -1
	}
	fmt.Printf("%s\n%s\n", cyan(endpoint), green(ballot))
	return ballot
}

func (node *Node) Campaign(ownPort int, term int, votes *int) {
	for i := 0; i < numNodes; i++ {
		//Ping all ports except self
		if minPort+i != ownPort {
			if node.requestVoteFromNode(fmt.Sprintf("%s:%v%s%s/%v/%v", url, (minPort+i), api, endRequestVote, term, ownPort)) == ownPort {
				*votes++
			}
		}
	}
}

// Wait for ping, if no ping received within timeout, promote self to candidate
func (node *Node) followerAction(ping <-chan int, rGen *rand.Rand) {
	select {
	case <-time.After(time.Duration(timeout+rGen.Intn(timeout)) * time.Millisecond): //Timeout
		fmt.Println("Promoted to Candidate")
		node.rank++
	case <-ping: //Pinged
		fmt.Println("Ping recieved")
	}
}

// Need to implement, automatically upgrade to leader for now
func (node *Node) candidateAction(ownPort int) {
	votes := 1 //Vote for self
	//increment term
	node.term++
	fmt.Printf("Campaign term: %v\n", yellow(node.term))
	//request votes from other nodes in go routine
	node.Campaign(ownPort, node.term, &votes)

	//Timeout
	//If achieved a simple majority, then promote self
	select {
	case <-time.After(time.Duration(timeout/2) * time.Millisecond):
		if votes >= numNodes/2+1 {
			node.rank++
		}
	}
}

// Ping all other nodes with empty appendLogs call periodically to prevent timeouts
func (node *Node) leaderAction(ownPort int) {
	select {
	case <-time.After(time.Duration(timeout/2) * time.Millisecond):
		node.HeartBeat(ownPort)
	}
}

func (node *Node) performRankAction(ping <-chan int, ownPort int, rGen *rand.Rand) {
	fmt.Println(yellow(node.rank.String()))
	switch node.rank {
	case Follower:
		node.followerAction(ping, rGen)
	case Candidate:
		node.candidateAction(ownPort)
	case Leader:
		node.leaderAction(ownPort)
	}
}

func main() {
	var ownPort int
	ping := make(chan int, 0)
	defer close(ping)

	node := Node{ping, 0, 0, 0}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		node.FindAndServePort(&ownPort)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rGen := rand.New(rand.NewSource(time.Now().UnixNano()))
		for {
			node.performRankAction(node.ping, ownPort, rGen)
		}
	}()

	wg.Wait()
}
