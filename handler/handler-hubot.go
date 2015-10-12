package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hico-horiuchi/sensu-plugins-go/lib/handler"
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

	request, err := http.NewRequest("POST", url(&h.Config), strings.NewReader(payload(&h.Event)))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")

	http.DefaultClient.Do(request)
}

func payload(event *handler.EventStruct) string {
	body, err := json.Marshal(metricsStruct{
		Client:      event.Client.Name,
		Check:       event.Check.Name,
		Output:      strings.TrimRight(event.Check.Output, "\n"),
		Status:      event.Check.Status,
		Occurrences: event.Occurrences,
	})
	if err != nil {
		return ""
	}

	return string(body)
}

func url(config *handler.ConfigStruct) string {
	host := config.GetPath("hubot", "host").MustString()
	port := config.GetPath("hubot", "port").MustInt()
	room := config.GetPath("hubot", "room").MustInt()

	return fmt.Sprintf("http://%s:%d/sensu?room=%d", host, port, room)
}
