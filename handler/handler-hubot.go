package main

import (
	"./sensu/plugin"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/http"
	"strings"
)

type metricsStruct struct {
	Client      string `json:"client"`
	Check       string `json:"check"`
	Output      string `json:"output"`
	Status      int    `json:"status"`
	Occurrences int    `json:"occurrences"`
}

func main() {
	handler := plugin.NewHandler("/etc/sensu/conf.d/handler-hubot.json")

	metrics := newMetrics(handler.Event)
	body, _ := json.Marshal(metrics)
	payload := strings.NewReader(string(body))

	url := createURL(handler.Config)
	request, _ := http.NewRequest("POST", url, payload)
	request.Header.Set("Content-Type", "application/json")

	http.DefaultClient.Do(request)
}

func newMetrics(event plugin.EventStruct) metricsStruct {
	metrics := metricsStruct{}

	metrics.Client = event.Client.Name
	metrics.Check = event.Check.Name
	metrics.Output = strings.TrimRight(event.Check.Output, "\n")
	metrics.Status = event.Check.Status
	metrics.Occurrences = event.Occurrences

	return metrics
}

func createURL(config simplejson.Json) string {
	host := config.GetPath("hubot", "host").MustString()
	port := config.GetPath("hubot", "port").MustInt()
	room := config.GetPath("hubot", "room").MustInt()

	return fmt.Sprintf("http://%s:%d/sensu?room=%d", host, port, room)
}
