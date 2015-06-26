package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/hico-horiuchi/ohgi/sensu"
)

type handlerStruct struct {
	Event  EventStruct
	Config ConfigStruct
}

type EventStruct struct {
	sensu.EventStruct
}

type ConfigStruct struct {
	simplejson.Json
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
