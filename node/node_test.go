package node

import (
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	mock_node "github.com/rijuCB/RAFT2/node/mocks"
	"github.com/stretchr/testify/require"
)

func TestRankEnum(t *testing.T) {
	rankTest := Follower
	require.Equal(t, rankTest.String(), "Follower")
	rankTest++
	require.Equal(t, rankTest.String(), "Candidate")
	rankTest++
	require.Equal(t, rankTest.String(), "Leader")
	rankTest++
	require.Equal(t, rankTest.String(), "unknown")
}

func TestNodeAPI(t *testing.T) {
	var (
		ctrl = gomock.NewController(t)
		cli  = mock_node.NewMockInode(ctrl)
	)
	ping := make(chan int, 1)
	defer close(ping)
	ping <- 1
	nodeTest := Node{Ping: ping, RandomGen: rand.New(rand.NewSource(0))}

	println(cli)

	//Follower test
	nodeTest.PerformRankAction()
	nodeTest.PerformRankAction()

	//Candidate tests
	nodeTest.PerformRankAction()

	//Leader tests
	nodeTest.Rank = Leader
	nodeTest.PerformRankAction()

}
