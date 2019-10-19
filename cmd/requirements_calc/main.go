package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/klaital/factorio-tools/recipe_lister"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	var recipeListerDirectory string
	var builderWhitelistPath string
	var targetRecipe string
	var builderCount float64
	var productivityPerSlot float64
	var speedMultiplierPerBuilder float64

	var searchRecipes string

	flag.StringVar(&searchRecipes, "search", "", "Search the recipe list for anything containing this word")
	flag.StringVar(&builderWhitelistPath, "builders", "builders.txt", "Limit builders to the ones named in this file")
	flag.StringVar(&recipeListerDirectory, "recipes", "recipe-lister", "Directory containing output from recipe-lister mod")
	flag.StringVar(&targetRecipe, "make", "", "Recipe to make")
	flag.Float64Var(&builderCount, "rate", 1.0, "Number of machines making it")
	flag.Float64Var(&productivityPerSlot, "productivity", 0.0, "Percent productivity per slot in each builder. Use decimal, e.g. 0.4 for +40% productivity per slot")
	flag.Float64Var(&speedMultiplierPerBuilder, "speed", 1.0, "Speed multiplier applied to every builder. Use 1.0 for 'no bonus', or 8.0 for '+800% bonus'")
	flag.Parse()

	// Read the AssemblingMachine data
	machines, err := recipe_lister.LoadBuildersFromDirectory(recipeListerDirectory)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to load Builders config")
		return
	}

	// Read the Recipes data
	recipes, err := recipe_lister.LoadRecipeFile(fmt.Sprintf("%s/recipe.json", recipeListerDirectory))
	if err != nil {
		logrus.WithError(err).Errorf("Failed to load Recipes config")
		return
	}

	config := CalcConfig{
		PreferredMachines: map[string]string{
			"advanced-crafting":   "assembling-machine-3",
			"centrifuging":        "centrifuge",
			"chemistry":           "chemical-plant",
			"crafting":            "assembling-machine-3",
			"crafting-with-fluid": "assembling-machine-3",
			"oil-processing":      "oil-refinery",
			"rocket-building":     "rocket-silo",
			"smelting":            "electric-furnace",
		},
		Recipes:                    recipes,
		Machines:                   machines,
		BuilderProductivityPerSlot: productivityPerSlot,   // Use a Productivity Module 3 in each builder's slot, if applicable
		BuilderSpeedBonus:          speedMultiplierPerBuilder, // Assume +800% speed from Beacons
	}

	// Do the search, if requested
	if len(searchRecipes) > 0 {
		fmt.Printf("Searching for %s among %d recipes...\n", searchRecipes, len(recipes))
		for _, recipe := range recipes {
			if strings.Contains(string(recipe.Name), searchRecipes) {
				fmt.Printf("%+v\n", recipe)
			}
		}
		return
	}

	if len(recipeListerDirectory) == 0 {
		logrus.Errorf("No directory specified for recipe-lister output\n")
		return
	}

	if len(targetRecipe) == 0 {
		logrus.Errorf("No recipe specified")
		return
	}

	//// Find the recipe
	//recipe, ok := recipes[targetRecipe]
	//if !ok {
	//	fmt.Printf("Recipe %s not found in set", targetRecipe)
	//	return
	//}
	//
	//// Find the machine
	//machine, ok := machines[recipe_lister.MachineName(specifiedMachine)]
	//if !ok {
	//	fmt.Printf("Machine %s not found in set", specifiedMachine)
	//	return
	//}

	//// Calculate number of machines to hit the target rate
	//cyclesPerSecond := machine.GetCraftingSpeed() / recipe.Energy
	//machineCount := targetRate / cyclesPerSecond
	//fmt.Printf("MachineCount: %f\n", machineCount)
	//
	//// Calculate Input/Output rates
	//inputsPerSecond := make(map[recipe_lister.ItemName]float64, 0)
	//outputsPerSecond := make(map[recipe_lister.ItemName]float64, 0)
	//
	//for _, input := range recipe.Ingredients {
	//	inputsPerSecond[input.Name] += float64(input.Amount) * cyclesPerSecond * machineCount
	//}
	//for _, output := range recipe.Products {
	//	outputsPerSecond[output.Name] += float64(output.Amount) * output.Probability * cyclesPerSecond * machineCount
	//}

	if _, err := os.Stat(builderWhitelistPath); err != nil {
		logrus.WithError(err).Debugf("No whitelist at %s", builderWhitelistPath)
	} else {
		whitelist, loadErr := LoadBuilderWhitelist(builderWhitelistPath)
		if loadErr != nil {
			logrus.Errorf("Failed to load whitelist: %v", loadErr)
			return
		}
		config.BuilderWhitelist = whitelist
	}

	inputs, outputs, builder, err := config.CalculateRates(builderCount, recipe_lister.RecipeName(targetRecipe))
	if err != nil {
		fmt.Printf("Failed to calculate rates: %s\n", err.Error())
		return
	}

	fmt.Printf("Producing %s in %f %s\n", targetRecipe, builderCount, builder.GetName())
	fmt.Print("\n         ----- Inputs: -----\n")
	for name, rate := range inputs {
		fmt.Printf("%25s\t%f\n", name, rate)
	}
	fmt.Print("\n         ----- Outputs: -----\n")
	for name, rate := range outputs {
		fmt.Printf("%25s\t%f\n", name, rate)
	}
}

type CalcConfig struct {
	PreferredMachines map[string]string
	Recipes           map[recipe_lister.RecipeName]recipe_lister.Recipe
	Machines          map[recipe_lister.MachineName]recipe_lister.Builder

	// Only consider builders from this (optional) whitelist
	BuilderWhitelist map[recipe_lister.MachineName]bool

	// Modules and Beacons support
	BuilderProductivityPerSlot float64 // What size productivity module to apply to builder slots. Set to 0 to disable Productivity modules
	BuilderSpeedBonus          float64 // Total speed bonus on each builder. Usually comes from Speed Modules in Beacons.
}

func LoadBuilderWhitelist(path string) (map[recipe_lister.MachineName]bool, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(fileBytes), "\n")
	builderSet := make(map[recipe_lister.MachineName]bool, 0)
	for _, line := range lines {
		machineName := strings.TrimSpace(line)
		builderSet[recipe_lister.MachineName(machineName)] = true
	}

	return builderSet, nil
}

func (config *CalcConfig) FindBestBuilder(name recipe_lister.RecipeName) (recipe_lister.Builder, error) {
	// Pull the actual recipe
	recipe, ok := config.Recipes[name]
	if !ok {
		return nil, errors.New("recipe not found in config set")
	}

	var bestBuilder recipe_lister.Builder
	consideredMachines := make([]recipe_lister.MachineName, 0)
	for i, builder := range config.Machines {
		if builder.SupportsCraftingCategory(recipe.CraftingCategory) {
			logrus.Debugf("Considered: %s", builder.GetName())
			consideredMachines = append(consideredMachines, builder.GetName())
			// If there is a populated whitelist, and this builder isn't on it, then skip it
			if len(config.BuilderWhitelist) > 0 && !config.BuilderWhitelist[builder.GetName()] {
				continue
			}
			if bestBuilder == nil {
				bestBuilder = config.Machines[i]
				continue
			}

			// Give priority to builders with more module support. Machines with more module
			// inventory are higher-level anyway.
			if builder.GetModuleInventoryCount() > bestBuilder.GetModuleInventoryCount() {
				bestBuilder = config.Machines[i]
				continue
			}
			if builder.GetCraftingSpeed() > bestBuilder.GetCraftingSpeed() {
				bestBuilder = config.Machines[i]
				continue
			}
		}
	}

	if bestBuilder == nil {
		return nil, fmt.Errorf("no builders out of %d considered support this recipe: %+v", len(consideredMachines), recipe)
	}

	return bestBuilder, nil
}

func (config *CalcConfig) CalculateRates(machineCount float64, targetRecipe recipe_lister.RecipeName) (inputs map[recipe_lister.ItemName]float64, outputs map[recipe_lister.ItemName]float64, builder recipe_lister.Builder, err error) {
	builder, err = config.FindBestBuilder(targetRecipe)
	if err != nil {
		return nil, nil, nil, err
	}
	recipe, ok := config.Recipes[targetRecipe]
	if !ok {
		return nil, nil, nil, fmt.Errorf("recipe not found: %s", targetRecipe)
	}

	// This calculates the input and output rates per builder
	rates := recipe.CalculateRates(builder, config.BuilderProductivityPerSlot, config.BuilderSpeedBonus)

	// Calculate the total I/O rates for all machines
	inputs = make(map[recipe_lister.ItemName]float64, len(recipe.Ingredients))
	outputs = make(map[recipe_lister.ItemName]float64, len(recipe.Products))

	for ingredientName, rate := range rates.Inputs {
		inputs[ingredientName] += rate * machineCount
	}
	for productName, rate := range rates.Outputs {
		outputs[productName] += rate * machineCount
	}
	return inputs, outputs, builder, nil
}

func (config *CalcConfig) FindBestRecipe(targetName recipe_lister.ItemName) (*recipe_lister.Recipe, error) {
	var bestRecipe recipe_lister.Recipe

	// Search for recipes that produce this item at the best speed
	for i, recipe := range config.Recipes {
		normalizedEnergy := recipe.NormalizedEnergyForProduct(targetName)
		if normalizedEnergy == math.MaxFloat64 {
			// Recipe does not produce this item
			continue
		}
		if bestRecipe.Energy == 0 || normalizedEnergy < bestRecipe.NormalizedEnergyForProduct(targetName) {
			bestRecipe = config.Recipes[i]
		}
	}

	if bestRecipe.NormalizedEnergyForProduct(targetName) == math.MaxFloat64 {
		return nil, errors.New("no recipe for target")
	}

	return &bestRecipe, nil
}

func (config *CalcConfig) CalculateRatesRecursive(machineCount float64, targetItem recipe_lister.ItemName) (inputs map[recipe_lister.ItemName]float64, outputs map[recipe_lister.ItemName]float64, builder recipe_lister.Builder, err error) {
	recipe, err := config.FindBestRecipe(targetItem)
	if err != nil {
		return nil, nil, nil, err
	}

	builder, err = config.FindBestBuilder(recipe.Name)
	if err != nil {
		return nil, nil, nil, err
	}

	cyclesPerSecond := builder.GetCraftingSpeed() / recipe.Energy
	inputs = make(map[recipe_lister.ItemName]float64, len(recipe.Ingredients))
	outputs = make(map[recipe_lister.ItemName]float64, len(recipe.Products))

	for _, ingredient := range recipe.Ingredients {
		inputs[ingredient.Name] += float64(ingredient.Amount) * cyclesPerSecond
	}
	for _, product := range recipe.Products {
		outputs[product.Name] += float64(product.Amount) * cyclesPerSecond * product.Probability
	}

	// TODO: Also find the total inputs/outputs for each ingredient
	//for _, ingredient := range recipe.Ingredients {
	//	ingredientInput, ingredientOutput, _, err := config.CalculateRatesRecursive()
	//}
	return inputs, outputs, builder, nil
}
