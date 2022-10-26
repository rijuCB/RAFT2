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

type IrestClient interface {
	SendEmptyAppendLogs(string)
	RequestVoteFromNode(string) int
}

type RestClient struct {
}

func (api *RestClient) SendEmptyAppendLogs(endpoint string) {
	resp, err := http.Post(endpoint, "application/json", strings.NewReader(""))

	if err != nil {
		fmt.Println(red(err))
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(red(err))
		return
	}
	fmt.Printf("%s\n%s\n", cyan(endpoint), green(string(b)))
}

func (api *RestClient) RequestVoteFromNode(endpoint string) int {
	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println(red(err))
		return -1
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(red(err))
		return -1
	}

	ballot, err := strconv.Atoi(string(b))
	if err != nil {
		fmt.Println(red(err))
		return -1
	}
	fmt.Printf("%s\n%s\n", cyan(endpoint), green(ballot))
	return ballot
}
