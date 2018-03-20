package factorio

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type CraftComponents struct {
	Items map[string]int
}

// CraftComponents returns a flat list of the items required to
// make the specified item directly.
func (db *ItemDb) ComputeCraftComponents(targetItemId string, quantity int, report *CraftComponents) (err error) {
	logger := log.WithFields(log.Fields{
		"func":   "ItemDb#ComputeCraftComponents",
		"ItemId": targetItemId,
	})
	item, ok := db.Data[targetItemId]
	if !ok {
		return errors.New(fmt.Sprintf("Item Not Found: %s", targetItemId))
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
		components[ingredient.ItemName] += ingredient.Quantity * quantity
	}

	// Initialize the report if needed. If not, merge our data in
	for itemId, qty := range components {
		logger.Debugln("Adding", itemId, " +", qty)
		report.Items[itemId] += qty
		// Compute the component breakdown for this item as well
		db.ComputeCraftComponents(itemId, qty, report)
	}

	return nil
}
