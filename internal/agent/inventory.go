package agent

import (
	"Neolithic/internal/core"
	"sort"
)

type Inventory []core.InventoryEntry

func (i *Inventory) GetAmount(res *core.Resource) int {
	for _, entry := range *i {
		if entry.Resource == res {
			return entry.Amount
		}
	}
	return 0
}

func (i *Inventory) AdjustAmount(res *core.Resource, amount int) {
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

	// need to add amount to Inventory
	if amount <= 0 {
		return // don't append negative amounts or zero
	}

	newEntry := core.InventoryEntry{Resource: res, Amount: amount}
	*i = append(*i, core.InventoryEntry{})
	copy((*i)[idx+1:], (*i)[idx:])
	(*i)[idx] = newEntry
}

func (i *Inventory) DeepCopy() core.Inventory {
	panic("implement me")
}

func (i *Inventory) String() string {
	panic("implement me")
}

func (i *Inventory) Entries() []core.InventoryEntry {
	copied := make([]core.InventoryEntry, len(*i))
	copy(copied, *i)
	return copied
}
