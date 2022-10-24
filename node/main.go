package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
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

func FindAndServePort() {
	// Colours!!
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	//create a new router
	router := mux.NewRouter()

	// specify endpoints, handler functions and HTTP method
	router.HandleFunc("/api/v1/append-entries", AppendLogs).Methods("POST")
	router.HandleFunc("/api/v1/request-vote", RequestVote).Methods("GET")
	http.Handle("/", router)

	var l net.Listener
	var err error
	//Search for open ports between 8091 & 8100
	for port := 8091; port < 8100; port++ {
		attemptPort := fmt.Sprintf(":%v", port)
		fmt.Printf("Attempting to open port %v ... ", cyan(port))
		l, err = net.Listen("tcp", attemptPort) //Attempt to open port
		if err == nil {
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

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		FindAndServePort()
	}()
	wg.Wait()
}
