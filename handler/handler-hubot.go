package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"../lib/handler"
)

type metricsStruct struct {
	Client      string `json:"client"`
	Check       string `json:"check"`
	Output      string `json:"output"`
	Status      int    `json:"status"`
	Occurrences int    `json:"occurrences"`
}

func main() {
	h := handler.New("/etc/sensu/conf.d/handler-hubot.json")

	request, err := http.NewRequest("POST", url(h.Config), strings.NewReader(payload(h.Event)))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")

	http.DefaultClient.Do(request)
}

func payload(event handler.EventStruct) string {
	metrics := metricsStruct{}

	metrics.Client = event.Client.Name
	metrics.Check = event.Check.Name
	metrics.Output = strings.TrimRight(event.Check.Output, "\n")
	metrics.Status = event.Check.Status
	metrics.Occurrences = event.Occurrences

	body, err := json.Marshal(metrics)
	if err != nil {
		return ""
	}

	return string(body)
}

func url(config handler.ConfigStruct) string {
	host := config.GetPath("hubot", "host").MustString()
	port := config.GetPath("hubot", "port").MustInt()
	room := config.GetPath("hubot", "room").MustInt()

	return fmt.Sprintf("http://%s:%d/sensu?room=%d", host, port, room)
}
