package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmptyAppendLogs(t *testing.T) {
	payload := ""
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, payload)
	}))
	defer svr.Close()

	client := RestClient{}

	client.SendEmptyAppendLogs(svr.URL)
}

func MockRequestServer(payload string) int {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, payload)
	}))
	defer svr.Close()

	client := RestClient{}

	return client.RequestVoteFromNode(svr.URL)
}

func TestRequestVoteFromNode(t *testing.T) {
	require.Equal(t, -1, MockRequestServer(""))
	require.Equal(t, 8092, MockRequestServer("8092"))
}
