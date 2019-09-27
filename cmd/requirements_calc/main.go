package main

import (
	"flag"
	"fmt"
	"github.com/klaital/factorio-tools/recipe_lister"
)

func main() {
	var recipeListerDirectory string
	var targetRecipe string
	var targetRate float64
	var specifiedMachine string

	flag.StringVar(&recipeListerDirectory, "recipes", "recipe-lister", "Directory containing output from recipe-lister mod")
	flag.StringVar(&targetRecipe, "make", "", "Recipe to make")
	flag.Float64Var(&targetRate, "rate", 1.0, "Rate to make them")
	flag.StringVar(&specifiedMachine, "machine", "assembling-machine-2", "Desired assembling machine")
	flag.Parse()

	if len(recipeListerDirectory) == 0 {
		fmt.Printf("No directory specified for recipe-lister output\n")
		return
	}

	if len(targetRecipe) == 0 {
		fmt.Printf("No recipe specified")
		return
	}

	// Read the AssemblingMachine data
	machines, err := recipe_lister.LoadAssemblingMachinesFile(fmt.Sprintf("%s/assembling-machine.json", recipeListerDirectory))
	if err != nil {
		fmt.Printf("Failed to load Machines config: %v", err)
		return
	}

	// Read the Recipes data
	recipes, err := recipe_lister.LoadRecipeFile(fmt.Sprintf("%s/recipe.json", recipeListerDirectory))
	if err != nil {
		fmt.Printf("Failed to load Recipes config: %v", err)
		return
	}

	// Find the recipe
	recipe, ok := recipes[targetRecipe]
	if !ok {
		fmt.Printf("Recipe %s not found in set", targetRecipe)
		return
	}

	// Find the machine
	machine, ok := machines[specifiedMachine]
	if !ok {
		fmt.Printf("Machine %s not found in set", specifiedMachine)
		return
	}

	// Calculate number of machines to hit the target rate
	cyclesPerSecond := machine.CraftingSpeed / recipe.Energy
	machineCount := targetRate / cyclesPerSecond
	fmt.Printf("MachineCount: %f\n", machineCount)
}
