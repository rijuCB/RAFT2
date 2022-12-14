package node

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	mock_client "github.com/rijuCB/RAFT2/restClient/mocks"
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

func TestNodeFollower(t *testing.T) {
	var (
		ctrl = gomock.NewController(t)
		cli  = mock_client.NewMockIrestClient(ctrl)
	)
	ping := make(chan int, 1)
	defer close(ping)
	nodeTest := Node{Ping: ping, OwnPort: 8091, RandomGen: rand.New(rand.NewSource(0)), API: cli}

	//Follower test
	//Remain follower after recieving ping
	ping <- 1
	nodeTest.PerformRankAction()
	require.Equal(t, nodeTest.Rank, Follower)

	//Upgrade to candidate if no ping received
	nodeTest.PerformRankAction()
	require.Equal(t, nodeTest.Rank, Candidate)

}

func TestNodeCandidate(t *testing.T) {
	var (
		ctrl = gomock.NewController(t)
		cli  = mock_client.NewMockIrestClient(ctrl)
	)
	ping := make(chan int, 1)
	defer close(ping)
	nodeTest := Node{Ping: ping, OwnPort: 8091, RandomGen: rand.New(rand.NewSource(0)), API: cli}

	//No votes - Remain candidate
	nodeTest.Rank = Candidate
	cli.EXPECT().RequestVoteFromNode(gomock.Any()).Return(-1, errors.New("err")).Times(2)
	nodeTest.PerformRankAction()
	require.Equal(t, nodeTest.Rank, Candidate)
	//1 vote - Become leader
	nodeTest.Rank = Candidate
	cli.EXPECT().RequestVoteFromNode(gomock.Any()).Return(-1, errors.New("err")).Times(1)
	cli.EXPECT().RequestVoteFromNode(gomock.Any()).Return(8091, nil).Times(1)
	nodeTest.PerformRankAction()
	require.Equal(t, nodeTest.Rank, Leader)
	//2 votes - Become leader
	nodeTest.Rank = Candidate
	cli.EXPECT().RequestVoteFromNode(gomock.Any()).Return(8091, nil).Times(2)
	nodeTest.PerformRankAction()
	require.Equal(t, nodeTest.Rank, Leader)

}

func TestNodeLeader(t *testing.T) {
	var (
		ctrl = gomock.NewController(t)
		cli  = mock_client.NewMockIrestClient(ctrl)
	)
	ping := make(chan int, 1)
	defer close(ping)
	nodeTest := Node{Ping: ping, OwnPort: 8091, RandomGen: rand.New(rand.NewSource(0)), API: cli}

	cli.EXPECT().SendEmptyAppendLogs(gomock.Any()).Times(2)
	nodeTest.Rank = Leader
	nodeTest.PerformRankAction()
}
