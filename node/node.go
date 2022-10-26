package node

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
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

// Constant vars
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

// Colours!!
var (
	cyan   = color.New(color.FgCyan).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

// Interfaces
//
//go:generate go run github.com/golang/mock/mockgen -destination mocks/Inode.go github.com/rijuCB/RAFT2/node Inode
type Inode interface {
	PerformRankAction()
	LeaderAction()
	CandidateAction()
	FollowerAction()
}

// Node struct
type Node struct {
	Ping chan int
	Rank Rank

	Term int //Add mutex
	Vote int

	OwnPort   int        //Stores ownPort addr
	RandomGen *rand.Rand //Random number gen
}

func (node *Node) AppendLogs(w http.ResponseWriter, r *http.Request) {
	if node.Rank == Follower {
		node.Ping <- 1
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

	if newTerm > node.Term {
		node.Vote = requester
		node.Term = newTerm
	}
	if node.Rank == Follower {
		node.Ping <- 1
	}

	fmt.Printf("Vote cast for %v\n", green(node.Vote))

	//update response
	fmt.Fprintf(w, "%v", node.Vote)
}

func (node *Node) FindAndServePort() {
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
			node.OwnPort = minPort + i
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

func (node *Node) sendEmptyAppendLogs(endpoint string) {
	resp, err := http.Post(endpoint, "application/json", strings.NewReader(""))

	if err != nil {
		fmt.Println(red(err))
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(red(err))
		return
	}
	fmt.Printf("%s\n%s\n", cyan(endpoint), green(string(b)))
}

func (node *Node) heartBeat() {
	for i := 0; i < numNodes; i++ {
		//Ping all ports except self
		if minPort+i != node.OwnPort {
			node.sendEmptyAppendLogs(fmt.Sprintf("%s:%v%s%s", url, (minPort + i), api, endAppendLogs))
		}
	}
}

func (node *Node) requestVoteFromNode(endpoint string) int {
	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println(red(err))
		return -1
	}
	defer resp.Body.Close()
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

func (node *Node) campaign(votes *int) {
	for i := 0; i < numNodes; i++ {
		//Ping all ports except self
		if minPort+i != node.OwnPort {
			if node.requestVoteFromNode(fmt.Sprintf("%s:%v%s%s/%v/%v", url, (minPort+i), api, endRequestVote, node.Term, node.OwnPort)) == node.OwnPort {
				*votes++
			}
		}
	}
}

// Wait for ping, if no ping received within timeout, promote self to candidate
func (node *Node) FollowerAction() {
	select {
	case <-time.After(time.Duration(timeout+node.RandomGen.Intn(timeout)) * time.Millisecond): //Timeout
		fmt.Println("Promoted to Candidate")
		node.Rank++
	case <-node.Ping: //Pinged
		fmt.Println("Ping recieved")
	}
}

// Need to implement, automatically upgrade to leader for now
func (node *Node) CandidateAction() {
	votes := 1 //Vote for self
	//increment term
	node.Term++
	fmt.Printf("Campaign term: %v\n", yellow(node.Term))
	//request votes from other nodes in go routine
	node.campaign(&votes)

	//Timeout
	//If achieved a simple majority, then promote self
	time.Sleep(time.Duration(timeout/2) * time.Millisecond)
	if votes >= numNodes/2+1 {
		node.Rank++
	}
}

// Ping all other nodes with empty appendLogs call periodically to prevent timeouts
func (node *Node) LeaderAction() {
	time.Sleep(time.Duration(timeout/2) * time.Millisecond)
	node.heartBeat()
}

func (node *Node) PerformRankAction() {
	fmt.Println(yellow(node.Rank.String()))
	switch node.Rank {
	case Follower:
		node.FollowerAction()
	case Candidate:
		node.CandidateAction()
	case Leader:
		node.LeaderAction()
	}
}
