package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/zwkno1/goflow"
)

type Hello struct {
	name string
}

func (h *Hello) Go(ctx context.Context) error {
	fmt.Println("Hello " + h.name + "!")
	return nil
}

type Factory struct {
}

func (f *Factory) Create(name string) (goflow.Node, error) {
	return &Hello{
		name: name,
	}, nil
}

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("read config error: " + err.Error())
		return
	}
	config := &goflow.Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		fmt.Println("unmarshal config error: " + err.Error())
		return
	}

	f, err := goflow.NewFlow(config, &Factory{})
	if err != nil {
		fmt.Println("create goflow error: " + err.Error())
		return
	}
	if err := f.Go(nil); err != nil {
		fmt.Println("go error: " + err.Error())
	}
}
