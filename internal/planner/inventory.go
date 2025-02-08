package planner

import "fmt"

// Inventory represents the amount of resources that an agent or location has
type Inventory map[*Resource]int

// Copy copies an inventory
func (i Inventory) Copy() Inventory {
	newInventory := make(Inventory)
	for k, v := range i {
		newInventory[k] = v
	}

	return newInventory
}

// String returns a string representation of an inventory.
func (i Inventory) String() string {
	output := ""
	for name, amount := range i {
		output += fmt.Sprintf("      %s: %d\n", name.Name, amount)
	}
	return output
}
