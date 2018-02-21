package factorio

import (
	"errors"
	"fmt"
)

type CraftComponents struct {
	Items	map[string]int
}

// CraftComponents returns a flat list of the items required to 
// make the specified item directly. No recursive computation 
// of those ingredients' ingredients.
func (db *ItemDb) CraftComponents(itemId string, report *CraftComponents) (err error) {
	item, ok := db.Data[itemId]
	if !ok {
		return errors.New(fmt.Sprintf("Item Not Found: %s", itemId))
	}

	// The item will have no recipes if it is a base component
	if len(item.Recipes) == 0 {
		return nil
	}

	components := make(map[string]int)
	// Only use the first recipe.
	for _, ingredient := range item.Recipes[0].Ingredients {
		// Utilizing the fact that a Go map access returns 
		// the zero-value of the result if the key is not found,
		// we don't need to check for membership and initialize.
		components[ingredient.ItemName] += ingredient.Quantity
	}

	// Initialize the report if needed. If not, merge our data in
	if report == nil {
		report = &CraftComponents{Items:components}
	} else {
		for itemId, qty := range components {
			report.Items[itemId] += qty
		}
	}

	return nil
}
