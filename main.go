package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/rijuCB/RAFT2/node"
)

func main() {
	ping := make(chan int, 0)
	defer close(ping)

	node := node.Node{Ping: ping, RandomGen: rand.New(rand.NewSource(time.Now().UnixNano()))}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		node.FindAndServePort()
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
