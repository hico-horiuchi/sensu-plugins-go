package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

type healthStruct struct {
	Status string
}

func main() {
	var (
		host    string
		port    int
		timeout int
	)

	pflag.StringVarP(&host, "host", "h", "localhost", "HOST")
	pflag.IntVarP(&port, "port", "p", 9200, "PORT")
	pflag.IntVarP(&timeout, "timeout", "t", 30, "TIMEOUT")
	pflag.Parse()

	url := "http://" + host + ":" + strconv.Itoa(port) + "/_cluster/health"
	status := getHealthStatus(url, timeout)

	switch status {
	case "green":
		fmt.Printf("CheckElasticsearch OK: Cluster is green\n")
		os.Exit(0)
	case "yellow":
		fmt.Printf("CheckElasticsearch WARNING: Cluster is yellow\n")
		os.Exit(1)
	case "red":
		fmt.Printf("CheckElasticsearch CRITICAL: Cluster is red\n")
		os.Exit(2)
	}
}

func getHealthStatus(url string, timeout int) string {
	var health healthStruct

	http.DefaultClient.Timeout = time.Duration(timeout) * time.Second
	request, _ := http.NewRequest("GET", url, nil)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		fmt.Println("CheckElasticsearch CRITICAL:", err)
		os.Exit(2)
	}
	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal(body, &health)
	return health.Status
}
