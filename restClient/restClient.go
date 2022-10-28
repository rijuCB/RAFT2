package rest

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	color "github.com/fatih/color"
)

// Colours!!
var (
	cyan   = color.New(color.FgCyan).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

//go:generate go run github.com/golang/mock/mockgen -destination mocks/IrestClient.go github.com/rijuCB/RAFT2/restClient IrestClient
type IrestClient interface {
	SendEmptyAppendLogs(string) error
	RequestVoteFromNode(string) (int, error)
}

type RestClient struct {
}

// Sends an empty AppendLogs request to reset the target nodes timeout
func (api *RestClient) SendEmptyAppendLogs(endpoint string) error {
	resp, err := http.Post(endpoint, "application/json", strings.NewReader(""))

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n%s\n", cyan(endpoint), green(string(b)))
	return nil
}

// Requests an external node for a vote
func (api *RestClient) RequestVoteFromNode(endpoint string) (int, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}
	//Ensure vote it a valid integer
	ballot, err := strconv.Atoi(string(b))
	if err != nil {
		return -1, err
	}
	fmt.Printf("%s\n%s\n", cyan(endpoint), green(ballot))
	return ballot, nil
}
