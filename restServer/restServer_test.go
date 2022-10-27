package rest

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mux "github.com/gorilla/mux"
	"github.com/rijuCB/RAFT2/node"
	"github.com/stretchr/testify/require"
)

func TestFindAndServePort(t *testing.T) {
	// ping := make(chan int, 1)
	// defer close(ping)
	// testNode := node.Node{Ping: ping}
	// testServer := RestServer{&testNode}

	// //Test FindAndServePort()
	// testServer.FindAndServePort()
}

func TestAppendLogs(t *testing.T) {
	ping := make(chan int, 1)
	defer close(ping)
	testNode := node.Node{Ping: ping}
	testServer := RestServer{&testNode}

	//Test AppendLogs(http.ResponseWriter, *http.Request)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	testServer.AppendLogs(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "API is up and running", string(data))

}

func SendMockRequest(testServer *RestServer, term string, requester string) ([]byte, error) {
	end := "http://localhost:8091/api/v1/request-vote/term/requester"
	param := map[string]string{
		"term":      "1",
		"requester": "8092",
	}

	req := httptest.NewRequest(http.MethodGet, end, nil)
	req = mux.SetURLVars(req, param)
	w := httptest.NewRecorder()
	testServer.RequestVote(w, req)
	res := w.Result()
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func TestRequestVote(t *testing.T) {
	ping := make(chan int, 1)
	defer close(ping)
	testNode := node.Node{Ping: ping, OwnPort: 8091}
	testServer := RestServer{&testNode}

	data, err := SendMockRequest(&testServer, "1", "8092")
	require.NoError(t, err)
	require.Equal(t, "8092", string(data))
}
