package main

import (
	"strings"

	"github.com/hico-horiuchi/sensu-plugins-go/lib/handler"
	"github.com/nlopes/slack"
)

func main() {
	h := handler.New("/etc/sensu/conf.d/handler-slack.json")
	api := slack.New(h.Config.GetPath("slack", "token").MustString())

	api.PostMessage(
		channelID(api, h.Config.GetPath("slack", "channel").MustString()),
		"",
		slack.PostMessageParameters{
			Username:    "Sensu",
			Attachments: []slack.Attachment{*attachment(&h.Event)},
			IconURL:     "https://sensuapp.org/img/sensu_flat_logo_large-ce32365a.png",
		},
	)
}

func channelID(api *slack.Client, name string) string {
	channels, err := api.GetChannels(true)
	if err != nil {
		return ""
	}

	name = strings.TrimLeft(name, "#")
	for _, channel := range channels {
		if channel.Name == name {
			return channel.ID
		}
	}

	return ""
}

func attachment(event *handler.EventStruct) *slack.Attachment {
	return &slack.Attachment{
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
	str = append(str, ("*Check* : " + event.Check.Name + "\n")...)
	str = append(str, ("```\n" + strings.TrimRight(event.Check.Output, "\n") + "\n```")...)

	return string(str)
}
