package rest

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rijuCB/RAFT2/node"
	"github.com/stretchr/testify/require"
)

func TestServerAPI(t *testing.T) {
	ping := make(chan int, 1)
	defer close(ping)
	testNode := node.Node{Ping: ping}
	testServer := RestServer{&testNode}

	//Test FindAndServePort()
	// testServer.FindAndServePort()

	//Test AppendLogs(http.ResponseWriter, *http.Request)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	testServer.AppendLogs(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, string(data), "API is up and running")

	//Test RequestVote(http.ResponseWriter, *http.Request)
}
