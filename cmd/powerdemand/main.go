package main

import (
	"flag"
	"fmt"
	"github.com/klaital/factorio-tools/factorio"
	"github.com/klaital/factorio-tools/recipe_lister"
	"io/ioutil"
)

func main() {
	var blueprintPath string
	var recipeListerDirectory string

	flag.StringVar(&blueprintPath, "bp", "", "File containing blueprint data")
	flag.StringVar(&recipeListerDirectory, "recipes", "", "Directory containing output from recipe-lister mod")
	flag.Parse()

	if len(blueprintPath) == 0 {
		fmt.Printf("No blueprint file given.\n")
		return
	}

	if len(recipeListerDirectory) == 0 {
		fmt.Printf("No directory specified for recipe-lister output\n")
		return
	}

	// Read the AssemblingMachine data
	machines, err := recipe_lister.LoadMachinesDirectory(recipeListerDirectory)
	if err != nil {
		fmt.Printf("Failed to load Machines config: %v", err)
		return
	}

	// Read the file
	bpBytes, err := ioutil.ReadFile(blueprintPath)
	if err != nil {
		fmt.Printf("Failed to read blueprint file: %v", err)
		return
	}
	blueprint, err := factorio.ParseBlueprintString(string(bpBytes))
	if err != nil {
		fmt.Printf("Failed to load BP from string: %v", err)
		return
	}

	// Enumerate and count the entities
	entities := make(map[string]int64)
	totalPowerForEntity := make(map[string]float64)
	totalPower := 0.0
	for _, entity := range blueprint.Details.Entities {
		entities[entity.Name] = entities[entity.Name] + 1
		if machine, ok := machines[entity.Name]; ok {
			totalPower += machine.GetOperatingKiloWatts()
			totalPowerForEntity[entity.Name] = totalPowerForEntity[entity.Name] + machine.GetOperatingKiloWatts()
		}
	}

	for entityName, count := range entities {
		if totalPowerForEntity[entityName] > 0.0 {
			fmt.Printf("%s\t%d\t%dkW\n", entityName, count, int64(totalPowerForEntity[entityName]))
		}
	}

	fmt.Printf("\nTotal Power:\t%fMW\n", totalPower/1000.0)
}

