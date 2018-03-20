package factorio

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Item struct {
	Name    string   `json:"name"`
	Recipes []Recipe `json:"recipes"`
}
type Recipe struct {
	Time        float64      `json:"time"`
	Yield       int          `json:"yield"`
	Ingredients []Ingredient `json:"ingredients"`
}
type Ingredient struct {
	Quantity int    `json:"qty"`
	ItemName string `json:"item"`
	//	Item		*Item
}

// The master data structure for manipulating the database in memory
type ItemDb struct {
	Data map[string]Item
}

func LoadJsonToDb(jsonDb string) (db *ItemDb, err error) {
	db = new(ItemDb)
	db.Data = make(map[string]Item)

	items := make([]Item, 0)
	err = json.Unmarshal([]byte(jsonDb), &items)

	// Add each item to the database
	for _, item := range items {
		log.Debugln("Adding record:", item.Name)
		db.Data[item.Name] = item
	}

	// add a pointer to each item used as an ingredient
	//	for _, item := range items {
	//		for _, recipe := range db.Data[item.Name].Recipes {
	//			for _, ingredient := range recipe.Ingredients {
	//				if tmpItem, ok := db.Data[ingredient.ItemName]; ok {
	//					ingredient.Item = db.Data[ingredient.ItemName]
	//				}
	//			}
	//		}
	//	}

	return db, err
}
