package goflow

type Factory interface {
	Create(name string) (Node, error)
}
