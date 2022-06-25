package goflow

type Config struct {
	Node     string   `json:"node,omitempty"`
	Parallel []Config `json:"parallel,omitempty"`
	Sequence []Config `json:"sequence,omitempty"`
}
