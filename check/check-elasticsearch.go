package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/hico-horiuchi/sensu-plugins-go/lib/check"
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

	c := check.New("CheckElasticsearch")
	c.Option.StringVarP(&host, "host", "h", "localhost", "HOST")
	c.Option.IntVarP(&port, "port", "P", 9200, "PORT")
	c.Option.IntVarP(&timeout, "timeout", "t", 30, "TIMEOUT")
	c.Init()

	status, err := healthStatus(host, port, timeout)
	if err != nil {
		c.Error(err)
	}

	switch status {
	case "green":
		c.Ok("Cluster is green")
	case "yellow":
		c.Warning("Cluster is yellow")
	case "red":
		c.Critical("Cluster is red")
	}
}

func healthStatus(host string, port int, timeout int) (string, error) {
	var health healthStruct
	http.DefaultClient.Timeout = time.Duration(timeout) * time.Second
	url := "http://" + host + ":" + strconv.Itoa(port) + "/_cluster/health"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	json.Unmarshal(body, &health)
	return health.Status, nil
}
