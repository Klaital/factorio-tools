package recipe_lister

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
)

type RecipeName string
type ItemName string
type Recipe struct {
	Name             RecipeName  `json:"name"`
	Energy           float64     `json:"energy"`
	Ingredients      []Component `json:"ingredients"`
	Products         []Component `json:"products"`
	CraftingCategory string      `json:"category"`
}
type Component struct {
	Type        string   `json:"type"`
	Name        ItemName `json:"name"`
	Amount      float64  `json:"amount"`
	Probability float64  `json:"probability"`
}

// NormalizedEnergyForProduct calculates the amount of energy required per each item produced
func (r *Recipe) NormalizedEnergyForProduct(productName ItemName) float64 {
	for _, product := range r.Products {
		if product.Name == productName {
			return r.Energy / float64(product.Amount) * product.Probability
		}
	}

	// Default return for "none found"
	return math.MaxFloat64
}

type RecipeRates struct {
	Inputs  map[ItemName]float64
	Outputs map[ItemName]float64
}

func (r *Recipe) CalculateRates(builder Builder, productivityPerSlot float64, speedMultiplier float64) RecipeRates {
	cyclesPerSecond := builder.GetCraftingSpeed() / r.Energy
	inputs := make(map[ItemName]float64)
	outputs := make(map[ItemName]float64)

	// Calculate actual productivity rate
	productivityRate := 1.0 + (float64(builder.GetModuleInventoryCount()) * productivityPerSlot)

	for _, ingredient := range r.Ingredients {
		inputs[ingredient.Name] = cyclesPerSecond * float64(ingredient.Amount) * speedMultiplier
	}
	for _, product := range r.Products {
		outputs[product.Name] = cyclesPerSecond * float64(product.Amount) * speedMultiplier * productivityRate * product.Probability
	}

	return RecipeRates{
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func LoadRecipeFile(path string) (dataSet map[RecipeName]Recipe, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening recipe file: %w")
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	recipeSet := make(map[RecipeName]Recipe, 0)
	if err = json.Unmarshal(b, &recipeSet); err != nil {
		return nil, fmt.Errorf("unmarshal recipe set: %w", err)
	}

	return recipeSet, nil
}
