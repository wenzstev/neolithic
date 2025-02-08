package planner

import (
	"fmt"
)

// Agent represents an entity that can do things in the world (i.e., a villager)
type Agent struct {
	Name string
}

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

// Resource represents a resource in the world
type Resource struct {
	Name string
}

// Location represents a location in the world
type Location struct {
	Name string
}

func TestAi() {
	wood := &Resource{
		Name: "wood",
	}
	axe := &Resource{
		Name: "axe",
	}
	forest := &Location{
		Name: "forest",
	}
	stockpile := &Location{
		Name: "stockpile",
	}
	woodcutter := &Agent{
		Name: "woodcutter",
	}

	start := &State{
		Locations: map[*Location]Inventory{
			forest:    {wood: 100},
			stockpile: {axe: 1},
		},
		Agents: map[*Agent]Inventory{
			woodcutter: {},
		},
	}

	goal := &State{
		Locations: map[*Location]Inventory{
			stockpile: {wood: 100},
		},
	}

	gatherWood := &Gather{
		resource: wood,
		amount:   1,
		location: forest,
		cost:     5,
	}

	gatherWoodWithAxe := &Gather{
		requires: axe,
		resource: wood,
		amount:   5,
		location: forest,
		cost:     1,
	}

	pickupAxe := &Gather{
		resource: axe,
		amount:   1,
		location: stockpile,
		cost:     1,
	}

	depositWood := &Deposit{
		resource: wood,
		amount:   50,
		location: stockpile,
	}

	actions := &Actions{
		gatherWood,
		gatherWoodWithAxe,
		pickupAxe,
		depositWood,
	}

	output, err := actions.AStar(start, goal, woodcutter, 100000000)
	if err != nil {
		panic(err)
	}
	fmt.Println("Found solution:")
	fmt.Println(output.String())
	fmt.Println("Final state: ")
	fmt.Println(output.expectedState.String())
}
