package core

import "fmt"

// Inventory is the minimal inventory interface that all packages need
type Inventory interface {
	fmt.Stringer
	GetAmount(res *Resource) int
	AdjustAmount(res *Resource, amount int)
	DeepCopy() Inventory
	Entries() []InventoryEntry
}

type InventoryEntry struct {
	Resource *Resource
	Amount   int
}
