package sensu

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
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
		fmt.Println("error reading event: ", err)
		os.Exit(1)
	}

	json.Unmarshal(bytes, &h.Event)
}

func (h *handlerStruct) loadConfig(path string) {
	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println("error loading config: ", err)
		os.Exit(1)
	}

	h.Config.UnmarshalJSON(bytes)
}
