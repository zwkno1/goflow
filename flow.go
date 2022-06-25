package goflow

import (
	"context"

	"github.com/pkg/errors"
)

type Flow struct {
	entrypoint Node
}

func NewFlow(config *Config, factory Factory) (*Flow, error) {
	node, err := createNode(config, factory)
	if err != nil {
		return nil, err
	}
	return &Flow{entrypoint: node}, nil
}

func (f *Flow) Go(ctx context.Context) error {
	return f.entrypoint.Go(ctx)
}

func createNode(config *Config, factory Factory) (Node, error) {
	switch {
	case len(config.Node) != 0:
		return factory.Create(config.Node)
	case len(config.Parallel) != 0 || len(config.Sequence) != 0:
		configs := config.Parallel
		if len(configs) == 0 {
			configs = config.Sequence
		}

		nodes := make([]Node, 0, len(configs))
		for _, c := range configs {
			node, err := createNode(&c, factory)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		}

		if len(config.Parallel) != 0 {
			return &Parallel{nodes: nodes}, nil
		} else {
			return &Sequence{nodes: nodes}, nil
		}
	default:
		return nil, errors.New("empty node")
	}
}
