package main

import (
	"flag"
	"fmt"
	"github.com/klaital/factorio-tools/recipe_lister"
	"os"
)

type CalculationConfig struct {
	Reactor   recipe_lister.Reactor
	Generator recipe_lister.Generator
	Boiler    recipe_lister.Boiler
}

type CalculationResults struct {
	ReactorQty    int
	GeneratorQty  int
	BoilerQty     int
	WaterRequired int
	MwYield       int64
	FuelConsumed  int64
}

func main() {
	var reactorCount int
	var recipeListerDirectory string
	var selectedMachines struct {
		boiler    string
		generator string
		reactor   string
	}

	flag.IntVar(&reactorCount, "count", 9, "Total number of desired reactors in a 3xN arrangement.")
	flag.StringVar(&recipeListerDirectory, "recipes", "recipe-lister", "Directory containing output from recipe-lister mod")
	flag.StringVar(&selectedMachines.reactor, "reactor", "fluid-reactor", "Which reactor to use. Pulled from recipe-lister/reactor.json")
	flag.StringVar(&selectedMachines.boiler, "boiler", "heat-exchanger", "Which boiler to use. Pulled from recipe-lister/boiler.json")
	flag.StringVar(&selectedMachines.generator, "generator", "steam-engine-2", "Which steam engine/turbine to use. Pulled from recipe-lister/generator.json")
	flag.Parse()

	if len(recipeListerDirectory) == 0 {
		fmt.Printf("No directory specified for recipe-lister output\n")
		return
	}

	// Load the machine data from recipelister
	boilers, err := recipe_lister.LoadBoilers(recipeListerDirectory)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	reactors, err := recipe_lister.LoadReactors(recipeListerDirectory)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	generators, err := recipe_lister.LoadGenerators(recipeListerDirectory)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	var calculator Calculator
	calculator = Calc3xN

	res := calculator(CalculationConfig{
		Reactor:   reactors[selectedMachines.reactor],
		Generator: generators[selectedMachines.generator],
		Boiler:    boilers[selectedMachines.boiler],
	}, reactorCount)

	fmt.Printf("%+v\n", res)
}
