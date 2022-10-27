package rest

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	mux "github.com/gorilla/mux"
	"github.com/rijuCB/RAFT2/node"
	"github.com/stretchr/testify/require"
)

func TestFindPortAndServePort(t *testing.T) {
	ping := make(chan int, 1)
	defer close(ping)
	testNode := node.Node{Ping: ping}
	testServer := RestServer{Node: &testNode}

	//Test FindPort()
	l := testServer.FindPort()
	require.Equal(t, "[::]:8091", l.Addr().String())

	l2 := testServer.FindPort()
	require.Equal(t, "[::]:8092", l2.Addr().String())

	l3 := testServer.FindPort()
	require.Equal(t, "[::]:8093", l3.Addr().String())

	ln := testServer.FindPort()
	require.Nil(t, ln)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		testServer.ServePort(l)
	}()
	time.Sleep(50) //Wait a bit for server to start

	//Check server is responsive
	endpoint := "http://localhost:8091/api/v1/request-vote/term/requester"
	_, err := http.Post(endpoint, "application/json", strings.NewReader(""))
	require.NoError(t, err)
	testServer.Server.Close()
	wg.Wait()
}

// Test AppendLogs(http.ResponseWriter, *http.Request)
func TestAppendLogs(t *testing.T) {
	ping := make(chan int, 1)
	defer close(ping)
	testNode := node.Node{Ping: ping}
	testServer := RestServer{Node: &testNode}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	testServer.AppendLogs(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "API is up and running", string(data))

}

func SendMockRequest(testServer *RestServer, term string, requester string) (string, error) {
	end := "http://localhost:8091/api/v1/request-vote/term/requester"
	param := map[string]string{
		"term":      term,
		"requester": requester,
	}

	req := httptest.NewRequest(http.MethodGet, end, nil)
	req = mux.SetURLVars(req, param)
	w := httptest.NewRecorder()
	testServer.RequestVote(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	return string(data), err
}

func TestRequestVote(t *testing.T) {
	ping := make(chan int, 1)
	defer close(ping)
	testNode := node.Node{Ping: ping, OwnPort: 8091}
	testServer := RestServer{Node: &testNode}

	var data string
	var err error

	//Test negative values
	data, err = SendMockRequest(&testServer, "-1", "8092")
	require.NoError(t, err)
	data, err = SendMockRequest(&testServer, "1", "-1")
	require.NoError(t, err)

	//Test invalid values
	data, err = SendMockRequest(&testServer, "seven", "8092")
	require.NoError(t, err)
	data, err = SendMockRequest(&testServer, "1", "eight")
	require.NoError(t, err)

	//Test base
	data, err = SendMockRequest(&testServer, "1", "8092")
	require.NoError(t, err)
	require.Equal(t, "8092", data)
	<-ping
	//Test already voted
	data, err = SendMockRequest(&testServer, "1", "8093")
	require.NoError(t, err)
	require.Equal(t, "8092", string(data))
	<-ping
	//Test new term
	data, err = SendMockRequest(&testServer, "2", "8093")
	println(testNode.Term)
	require.NoError(t, err)
	require.Equal(t, "8093", string(data))
	<-ping
}
