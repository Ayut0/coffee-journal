package domain

import "fmt"


type Process string

const (
	ProcessWashed Process = "Washed"
	ProcessNatural Process = "Natural"
	ProcessHoney Process = "Honey"
)

func NewProcess(s string) (Process, error) {
	switch Process(s) {
		case ProcessWashed, ProcessNatural, ProcessHoney:
			return Process(s), nil
		default:
			return "", fmt.Errorf("invalid process %q: must be Washed, Natural, or Honey", s)
	}
}

func (p Process) String() string {
	return string(p)
}