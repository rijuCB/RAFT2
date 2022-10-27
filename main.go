package main

import (
	"math/rand"
	"sync"
	"time"

	node "github.com/rijuCB/RAFT2/node"
	restClient "github.com/rijuCB/RAFT2/restClient"
	restServer "github.com/rijuCB/RAFT2/restServer"
)

func main() {
	ping := make(chan int, 0)
	defer close(ping)

	node := node.Node{Ping: ping, RandomGen: rand.New(rand.NewSource(time.Now().UnixNano())), API: &restClient.RestClient{}}
	restAPI := restServer.RestServer{Node: &node}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		restAPI.ServePort(restAPI.FindPort())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			node.PerformRankAction()
		}
	}()

	wg.Wait()
}
