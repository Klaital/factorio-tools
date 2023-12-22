package main

import (
	"flag"
	"fmt"
	"github.com/klaital/factorio-tools/recipe_lister"
	"os"
)

func main() {
	var err error

	var recipeListerDirectory string
	var machineId string
	var recipeId string
	var machineCount float64
	var listFile string

	flag.StringVar(&recipeListerDirectory, "recipes", "recipe-lister", "Directory containing output from recipe-lister mod")
	flag.StringVar(&machineId, "machine", "", "ID of the machine to use")
	flag.StringVar(&recipeId, "recipe", "", "ID of the recipe to implement")
	flag.Float64Var(&machineCount, "count", 1, "Number of machines to run")
	flag.StringVar(&listFile, "file", "", "Load processes from a file")
	flag.Parse()

	data, err := recipe_lister.LoadAll(recipeListerDirectory)
	if err != nil {
		panic(err)
	}

	var processes *recipe_lister.ProcessChain
	if len(listFile) > 0 {
		fmt.Printf("Loading processes from file %s...\n", listFile)
		processes, err = recipe_lister.LoadProcessChain(listFile, recipeListerDirectory)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Loaded %d processes\n", len(processes.Processes))
	} else {

		if len(machineId) == 0 {
			fmt.Printf("Must specify a machine ID")
			os.Exit(1)
		}
		if len(recipeId) == 0 {
			fmt.Printf("Must specify a recipe ID")
			os.Exit(1)
		}
		recipe, ok := data.Recipes[recipe_lister.RecipeName(recipeId)]
		if !ok {
			fmt.Printf("Invalid recipe ID")
			os.Exit(1)
		}
		machine, ok := data.Machines[machineId]
		if !ok {
			fmt.Printf("Invalid machine ID")
			os.Exit(1)
		}

		processes = &recipe_lister.ProcessChain{Processes: []recipe_lister.Process{
			{
				MachineCount: machineCount,
				Machine:      machine,
				Recipe:       recipe,
			},
		}}
	}

	overallRates := processes.TotalIO()
	fmt.Printf("==== Overall I/O ====\n")
	fmt.Printf("---- Inputs ----\n")
	for item, rate := range overallRates.Inputs {
		fmt.Printf("%s\t%f /s\n", item, rate)
	}
	fmt.Printf("---- Outputs ----\n")
	for item, rate := range overallRates.Outputs {
		fmt.Printf("%s\t%f /s\n", item, rate)
	}

	// TODO: display per-process I/O
	
	//for i, rates := range processes.() {
	//	comment := ""
	//	if len(processes.Processes[i].Notes) > 0 {
	//		comment = fmt.Sprintf(" (%s) ", processes.Processes[i].Notes)
	//	}
	//	fmt.Printf("\n===== %f x %s%s===== \n", processes.Processes[i].Count, processes.Processes[i].RecipeId, comment)
	//	fmt.Printf("---- Input\n")
	//	for item, rate := range rates.Inputs {
	//		fmt.Printf("%s\t%f /s\n", item, rate*machineCount)
	//	}
	//	fmt.Printf("---- Output\n")
	//	for item, rate := range rates.Outputs {
	//		fmt.Printf("%s\t%f /s\n", item, rate*machineCount)
	//	}
	//}
}
