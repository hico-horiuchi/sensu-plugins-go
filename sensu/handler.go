package sensu

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"os"
)

type handlerStruct struct {
	Event  EventStruct
	Config simplejson.Json
}

func NewHandler(configPath string) *handlerStruct {
	handler := &handlerStruct{}
	handler.readEvent()
	handler.loadConfig(configPath)
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
