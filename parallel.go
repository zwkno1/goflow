package goflow

import (
	"context"
	"sync"
)

type Parallel struct {
	nodes []Node
}

func (p *Parallel) Go(ctx context.Context) error {
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
