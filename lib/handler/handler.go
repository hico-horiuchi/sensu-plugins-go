package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	simplejson "github.com/bitly/go-simplejson"
)

type handlerStruct struct {
	Event  EventStruct
	Config ConfigStruct
}

type EventStruct struct {
	Client      clientStruct
	Check       checkStruct
	Occurrences int
	Action      string
}

type ConfigStruct struct {
	simplejson.Json
}

type clientStruct struct {
	Name          string
	Address       string
	Subscriptions []string
	Timestamp     int64
}

type checkStruct struct {
	Name        string
	Issued      int
	Output      string
	Status      int
	Command     string
	Subscribers []string
	Handler     string
	History     []string
	Flapping    bool
}

func New(path string) *handlerStruct {
	handler := &handlerStruct{}
	handler.readEvent()
	handler.loadConfig(path)

	return handler
}

func (h *handlerStruct) readEvent() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(bytes, &h.Event)
}

func (h *handlerStruct) loadConfig(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	h.Config.UnmarshalJSON(bytes)
}
