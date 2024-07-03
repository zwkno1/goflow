package goflow

import (
	"sync"

	"github.com/pkg/errors"
)

type Node[Context any] interface {
	Go(ctx Context) error
}

type parallel[Context any] struct {
	nodes []Node[Context]
}

func (p *parallel[Context]) Go(ctx Context) error {
	var wg sync.WaitGroup
	wg.Add(len(p.nodes))

	errs := make([]error, len(p.nodes))

	for i := range p.nodes {
		err := &errs[i]
		node := p.nodes[i]
		go func() {
			defer wg.Done()
			*err = node.Go(ctx)
		}()
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

type sequence[Context any] struct {
	nodes []Node[Context]
}

func (s *sequence[Context]) Go(ctx Context) error {
	for _, node := range s.nodes {
		if err := node.Go(ctx); err != nil {
			return err
		}
	}
	return nil
}

type NodeFactory[Context any] interface {
	Create(name string) (Node[Context], error)
}

type Flow[Context any] struct {
	entrypoint Node[Context]
}

type Config struct {
	Node     string   `json:"node,omitempty"`
	Parallel []Config `json:"parallel,omitempty"`
	Sequence []Config `json:"sequence,omitempty"`
}

func NewFlow[Context any](config *Config, factory NodeFactory[Context]) (*Flow[Context], error) {
	node, err := createNode(config, factory)
	if err != nil {
		return nil, err
	}
	return &Flow[Context]{entrypoint: node}, nil
}

func (f *Flow[Context]) Go(ctx Context) error {
	return f.entrypoint.Go(ctx)
}

func createNode[Context any](config *Config, factory NodeFactory[Context]) (Node[Context], error) {
	switch {
	case len(config.Node) != 0:
		return factory.Create(config.Node)
	case len(config.Parallel) != 0 || len(config.Sequence) != 0:
		configs := config.Parallel
		if len(configs) == 0 {
			configs = config.Sequence
		}

		nodes := make([]Node[Context], 0, len(configs))
		for _, c := range configs {
			node, err := createNode(&c, factory)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		}

		if len(config.Parallel) != 0 {
			return &parallel[Context]{nodes: nodes}, nil
		} else {
			return &sequence[Context]{nodes: nodes}, nil
		}
	default:
		return nil, errors.New("empty node")
	}
}
