package core

import "fmt"

type Agent interface {
	fmt.Stringer
	Name() string
	DeepCopy() Agent
	Inventory() Inventory
}
