package recipe_lister

import (
	"encoding/json"
	"io/ioutil"
)

type Recipe struct {
	Name string `json:"name"`
	Energy float64 `json:"energy"`
	Ingredients []Component `json:"ingredients"`
	Products []Component `json:"products"`
}

type Component struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Amount int64 `json:"amount"`
	Probability float64 `json:"probability"`
}

func LoadRecipeFile(path string) (dataSet map[string]Recipe, err error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	recipeSet := make(map[string]Recipe, 0)
	err = json.Unmarshal(fileBytes, &recipeSet)
	if err != nil {
		return nil, err
	}

	return recipeSet, nil
}