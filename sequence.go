package goflow

import "context"

type Sequence struct {
	nodes []Node
}

func (s *Sequence) Go(ctx context.Context) error {
	for _, node := range s.nodes {
		if err := node.Go(ctx); err != nil {
			return err
		}
	}
	return nil
}
