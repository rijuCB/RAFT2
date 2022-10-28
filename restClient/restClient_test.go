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

func MockRequestServer(payload string) (int, error) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, payload)
	}))
	defer svr.Close()

	client := RestClient{}

	return client.RequestVoteFromNode(svr.URL)
}

func TestRequestVoteFromNode(t *testing.T) {
	var vote int
	var err error
	vote, err = MockRequestServer("")
	require.Equal(t, -1, vote)
	require.Error(t, err)
	vote, err = MockRequestServer("8092")
	require.Equal(t, 8092, vote)
	require.NoError(t, err)
}
