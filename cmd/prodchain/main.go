package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/klaital/factorio-tools/recipe_lister"
	"os"
)

func main() {
	var recipeListerDirectory string
	var processesFile string
	flag.StringVar(&recipeListerDirectory, "recipes", "recipe-lister", "Directory containing output from recipe-lister mod")
	flag.StringVar(&processesFile, "processes", "processes.yml", "Config file containing the list of processes to run.")
	flag.Parse()

	gameData, err := recipe_lister.LoadAll(recipeListerDirectory)
	if err != nil {
		fmt.Printf("Failed to load game data: %+v", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded game data. %d machines, %d recipes\n", len(gameData.Machines), len(gameData.Recipes))

	processData, err := LoadProcessList(processesFile)
	if err != nil {
		fmt.Printf("Failed to load process data: %+v", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded %d proceses\n", len(processData.Processes))
}

type ProcessList struct {
	Processes []struct {
		ProcessId string
		Qty       float64
		MachineId string
	}
}

func LoadProcessList(path string) (*ProcessList, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading process file: %w", err)
	}
	var p ProcessList
	if err = json.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("parsing process file: %w", err)
	}
	return &p, nil
}
