package check

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type checkStruct struct {
	Name   string
	Option *pflag.FlagSet
}

func New(name string) *checkStruct {
	check := &checkStruct{
		Name:   name,
		Option: pflag.NewFlagSet(name, 1),
	}

	return check
}

func (c checkStruct) Init() {
	c.Option.Parse(os.Args[1:])
}

func (c checkStruct) Ok(output string) {
	fmt.Println(c.Name, "OK:", output)
	os.Exit(0)
}

func (c checkStruct) Warning(output string) {
	fmt.Println(c.Name, "WARNING:", output)
	os.Exit(1)
}

func (c checkStruct) Critical(output string) {
	fmt.Println(c.Name, "CRITICAL:", output)
	os.Exit(2)
}

func (c checkStruct) Error(err error) {
	fmt.Println(c.Name, "ERROR:", err)
	os.Exit(3)
}
