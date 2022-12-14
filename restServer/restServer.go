package rest

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	color "github.com/fatih/color"
	mux "github.com/gorilla/mux"
	node "github.com/rijuCB/RAFT2/node"
)

// Constant vars
const (
	minPort  = 8091
	numNodes = 3
	timeout  = 1500

	// url              = "http://localhost"
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

//go:generate go run github.com/golang/mock/mockgen -destination mocks/IRestServer.go github.com/rijuCB/RAFT2/restServer IRestServer
type IRestServer interface {
	FindAndServePort()
	AppendLogs(http.ResponseWriter, *http.Request)
	RequestVote(http.ResponseWriter, *http.Request)
}

type RestServer struct {
	Node   *node.Node
	Server *http.Server //Add this to allow closing the server externally
}

// Append logs to node (Not implemented)
// Resets node timeout
func (api *RestServer) AppendLogs(w http.ResponseWriter, r *http.Request) {
	if api.Node.Rank == node.Follower {
		api.Node.Ping <- 1
	}

	//specify status code
	w.WriteHeader(http.StatusOK)

	//update response writer
	fmt.Fprintf(w, "API is up and running")
}

// Parse http request and confirm they are integer values
func (api *RestServer) parseVoteRequest(r *http.Request) (int, int) {
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

// Responds to a vote request
// If first candidate in a term, votes for them
// Voted candidate remains consistent through a term
// Also resets node timeout
func (api *RestServer) RequestVote(w http.ResponseWriter, r *http.Request) {
	newTerm, requester := api.parseVoteRequest(r)
	if newTerm < 0 || requester < 0 {
		fmt.Fprintf(w, red("Invalid parameters!"))
		return //Invalid values
	}

	//specify HTTP status code
	w.WriteHeader(http.StatusOK)

	fmt.Printf("Received vote request from node:%v for term %v\n", cyan(requester), yellow(newTerm))

	if newTerm > api.Node.Term {
		api.Node.Vote = requester
		api.Node.Term = newTerm
	}
	if api.Node.Rank == node.Follower {
		api.Node.Ping <- 1
	}

	fmt.Printf("Vote cast for %v\n", green(api.Node.Vote))

	//update response
	fmt.Fprintf(w, "%v", api.Node.Vote)
}

// Search for an open port and return a listener with the reserved port
func (api *RestServer) FindPort() net.Listener {
	var l net.Listener = nil
	var err error
	//Search for open ports
	for i := 0; i < numNodes; i++ {
		attemptPort := fmt.Sprintf(":%v", (minPort + i))
		fmt.Printf("Attempting to open port %v ... ", cyan(attemptPort))
		l, err = net.Listen("tcp", attemptPort) //Attempt to open port
		if err == nil {
			api.Node.OwnPort = minPort + i
			fmt.Printf(green("success!\n"))
			break
		}
		fmt.Printf("%v\n%v\n", red("failed!"), yellow(err))
	}

	if l == nil {
		fmt.Printf(red("Unable to open port"))
		return nil
	}

	return l
}

// Create a REST server on a provided reserved port that API calls can be made to
func (api *RestServer) ServePort(l net.Listener) {
	//create a new router
	router := mux.NewRouter()

	// specify endpoints, handler functions and HTTP method
	router.HandleFunc(apiURL+endAppendLogs+paramAppendLogs, api.AppendLogs).Methods("POST")
	router.HandleFunc(apiURL+endRequestVote+paramRequestVote, api.RequestVote).Methods("GET")
	http.Handle("/", router)

	api.Server = &http.Server{Handler: router}

	err := api.Server.Serve(l)
	if err != nil { //Respond to requests
		fmt.Printf("ERROR!\n%v\n", red(err))
	}

}
