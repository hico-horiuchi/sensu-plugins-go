package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hico-horiuchi/sensu-plugins-go/lib/handler"
)

type attachmentStruct struct {
	Color      string   `json:"color"`
	Fallback   string   `json:"fallback"`
	Text       string   `json:"text"`
	MarkdownIn []string `json:"mrkdwn_in"`
}

type payloadStruct struct {
	Username    string             `json:"username"`
	Attachments []attachmentStruct `json:"attachments"`
	IconURL     string             `json:"icon_url"`
}

func main() {
	h := handler.New("/etc/sensu/conf.d/handler-slack.json")

	request, err := http.NewRequest(
		"POST",
		h.Config.GetPath("slack", "webhook_url").MustString(),
		strings.NewReader(payload(&h.Event)),
	)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")

	http.DefaultClient.Do(request)
}

func payload(event *handler.EventStruct) string {
	body, err := json.Marshal(payloadStruct{
		Username:    "Sensu",
		Attachments: []attachmentStruct{*attachment(event)},
		IconURL:     "https://sensuapp.org/img/sensu_flat_logo_large-ce32365a.png",
	})
	if err != nil {
		return ""
	}

	return string(body)
}

func attachment(event *handler.EventStruct) *attachmentStruct {
	return &attachmentStruct{
		Color:      color(event.Check.Status),
		Fallback:   event.Check.Name + " - " + event.Client.Name + " (" + strings.TrimRight(event.Check.Output, "\n") + ")",
		Text:       text(event),
		MarkdownIn: []string{"text"},
	}
}

func color(status int) string {
	switch status {
	case 0:
		return "#43ac6a"
	case 1:
		return "#f9ba46"
	case 2:
		return "#ea5443"
	}

	return "#9c9990"
}

func text(event *handler.EventStruct) string {
	var str []byte

	str = append(str, ("*Client* : " + event.Client.Name + "\n")...)
	str = append(str, ("*Address* : " + event.Client.Address + "\n")...)
	str = append(str, ("*Subscriptions* : " + strings.Join(event.Client.Subscriptions, ", ") + "\n")...)
	str = append(str, ("*Check* : " + event.Check.Name + "\n")...)
	str = append(str, ("```\n" + strings.TrimRight(event.Check.Output, "\n") + "\n```")...)

	return string(str)
}
