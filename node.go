package goflow

import "context"

type Node interface {
	Go(ctx context.Context) error
}
