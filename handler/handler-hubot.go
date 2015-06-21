package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"../sensu-plugin/handler"
)

type metricsStruct struct {
	Client      string `json:"client"`
	Check       string `json:"check"`
	Output      string `json:"output"`
	Status      int    `json:"status"`
	Occurrences int    `json:"occurrences"`
}

func main() {
	h := handler.New("/etc/sensu/conf.d/h-hubot.json")

	metrics := newMetrics(h.Event)
	body, _ := json.Marshal(metrics)
	payload := strings.NewReader(string(body))

	url := createURL(h.Config)
	request, _ := http.NewRequest("POST", url, payload)
	request.Header.Set("Content-Type", "application/json")

	http.DefaultClient.Do(request)
}

func newMetrics(event handler.EventStruct) metricsStruct {
	metrics := metricsStruct{}

	metrics.Client = event.Client.Name
	metrics.Check = event.Check.Name
	metrics.Output = strings.TrimRight(event.Check.Output, "\n")
	metrics.Status = event.Check.Status
	metrics.Occurrences = event.Occurrences

	return metrics
}

func createURL(config handler.ConfigStruct) string {
	host := config.GetPath("hubot", "host").MustString()
	port := config.GetPath("hubot", "port").MustInt()
	room := config.GetPath("hubot", "room").MustInt()

	return fmt.Sprintf("http://%s:%d/sensu?room=%d", host, port, room)
}
