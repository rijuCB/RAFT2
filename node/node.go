package node

import (
	"fmt"
	"math/rand"
	"time"

	color "github.com/fatih/color"
	restClient "github.com/rijuCB/RAFT2/restClient"
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
	apiURL           = "/api/v1"
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
	API       restClient.IrestClient
}

func (node *Node) heartBeat() {
	for i := 0; i < numNodes; i++ {
		//Ping all ports except self
		if minPort+i != node.OwnPort {
			node.API.SendEmptyAppendLogs(fmt.Sprintf("%s:%v%s%s", url, (minPort + i), apiURL, endAppendLogs))
		}
	}
}

func (node *Node) campaign(votes *int) {
	for i := 0; i < numNodes; i++ {
		//Ping all ports except self
		if minPort+i != node.OwnPort {
			if node.API.RequestVoteFromNode(fmt.Sprintf("%s:%v%s%s/%v/%v", url, (minPort+i), apiURL, endRequestVote, node.Term, node.OwnPort)) == node.OwnPort {
				*votes++
			}
		}
	}
}

// Wait for ping, if no ping received within timeout, promote self to candidate
func (node *Node) FollowerAction() {
	select {
	case <-time.After(time.Duration(timeout+node.RandomGen.Intn(timeout)) * time.Millisecond): //Timeout
		fmt.Println(red("Ping not recieved"))
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
