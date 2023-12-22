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
	AmountMin   float64  `json:"amount_min"`
	AmountMax   float64  `json:"amount_max"`
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

func NewRates() RecipeRates {
	return RecipeRates{
		Inputs:  map[ItemName]float64{},
		Outputs: map[ItemName]float64{},
	}
}

func (r *RecipeRates) ModifyInputs(lambda func(x float64) float64) {
	for k, v := range r.Inputs {
		r.Inputs[k] = lambda(v)
	}
}
func (r *RecipeRates) ModifyOutputs(lambda func(x float64) float64) {
	for k, v := range r.Outputs {
		r.Outputs[k] = lambda(v)
	}
}
func (r *RecipeRates) ModifyAll(lambda func(x float64) float64) {
	r.ModifyInputs(lambda)
	r.ModifyOutputs(lambda)
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
		basePerCycle := product.Amount * product.Probability
		if product.Amount <= 0 {
			basePerCycle = (product.AmountMin + product.AmountMax) * product.Probability
		}
		outputs[product.Name] = cyclesPerSecond * basePerCycle * speedMultiplier * productivityRate
	}

	return RecipeRates{
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func LoadRecipeFile(path string) (map[RecipeName]Recipe, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening recipe file: %w", err)
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

func LoadRecipes(directory string) (map[RecipeName]Recipe, error) {
	return LoadRecipeFile(fmt.Sprintf("%s/recipe.json", directory))
}

type GameData struct {
	Recipes  map[RecipeName]Recipe
	Machines map[string]AssemblingMachine
}

func LoadAll(directory string) (*GameData, error) {
	var err error
	var resp GameData

	if resp.Recipes, err = LoadRecipes(directory); err != nil {
		return nil, err
	}
	if resp.Machines, err = LoadAssemblingMachinesFile(fmt.Sprintf("%s/assembling-machine.json", directory)); err != nil {
		return nil, err
	}
	return &resp, nil
}
