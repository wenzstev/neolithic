package core

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// InventoryEntry represents a single resource and its quantity in an inventory.
type InventoryEntry struct {
	Resource *Resource
	Amount   int
}

// Inventory defines an interface for managing a collection of resources with their quantities.
// It provides methods for retrieving, adjusting, and copying inventory contents. Inventory contents are expected
// to be sorted by resource name.
type Inventory interface {
	fmt.Stringer
	// GetAmount retrieves the quantity of a specific resource from the inventory.
	GetAmount(res *Resource) int
	// AdjustAmount modifies the quantity of a specific resource in the inventory. Amount can be positive or negative.
	AdjustAmount(res *Resource, amount int)
	// DeepCopy creates a deep copy of the inventory. Resources are NOT deep copied,  but quantities are.
	DeepCopy() Inventory
	// Entries returns a slice of InventoryEntry, sorted by resource name. This is a copy of the inventory, not a reference.
	Entries() []InventoryEntry
}

// NewInventory creates and returns a new empty Inventory instance.
func NewInventory() Inventory {
	return &inventory{}
}

// inventory is an implementation of the Inventory interface.
type inventory []InventoryEntry

// GetAmount retrieves the quantity of a specific resource from the inventory.
func (i *inventory) GetAmount(res *Resource) int {
	for _, entry := range *i {
		if entry.Resource == res {
			return entry.Amount
		}
	}
	return 0
}

// AdjustAmount modifies the quantity of a specific resource in the inventory. Amount can be positive or negative.
func (i *inventory) AdjustAmount(res *Resource, amount int) {
	idx := sort.Search(len(*i), func(j int) bool {
		return (*i)[j].Resource.Name >= res.Name
	})

	if idx < len(*i) && (*i)[idx].Resource.Name == res.Name {
		(*i)[idx].Amount += amount
		if (*i)[idx].Amount <= 0 {
			*i = append((*i)[:idx], (*i)[idx+1:]...)
		}
		return
	}

	// need to add amount to inventory
	if amount <= 0 {
		return // don't append negative amounts or zero
	}

	newEntry := InventoryEntry{Resource: res, Amount: amount}
	*i = append(*i, InventoryEntry{})
	copy((*i)[idx+1:], (*i)[idx:])
	(*i)[idx] = newEntry
}

// DeepCopy creates a deep copy of the inventory. Resources are NOT deep copied.
func (i *inventory) DeepCopy() Inventory {
	copyInv := make(inventory, len(*i))
	for idx := 0; idx < len(*i); idx++ {
		copyInv[idx] = InventoryEntry{Resource: (*i)[idx].Resource, Amount: (*i)[idx].Amount}
	}
	return &copyInv
}

// String returns a string representation of the inventory.
func (i *inventory) String() string {
	if len(*i) == 0 {
		return "{}"
	}

	var sb strings.Builder

	for _, entry := range *i {
		// Use Fprintf to write directly to the builder's underlying buffer.
		// This avoids intermediate string allocations for each line.
		sb.WriteString("  ")
		sb.WriteString(entry.Resource.Name)
		sb.WriteString(": ")
		sb.WriteString(strconv.Itoa(entry.Amount))
	}
	return sb.String() // Creates the final string once from the buffer
}

// Entries returns a slice of InventoryEntry, sorted by resource name. This is a copy of the inventory, not a reference.
func (i *inventory) Entries() []InventoryEntry {
	copied := make([]InventoryEntry, len(*i))
	copy(copied, *i)
	return copied
}
