package main

import (
	"fmt"
	"github.com/klaital/factorio-tools/recipe_lister"
	"gopkg.in/yaml.v3"
	"os"
)

type ProcessDescription struct {
	Notes     string                   `yaml:"Notes"`
	MachineId string                   `yaml:"MachineId"`
	RecipeId  recipe_lister.RecipeName `yaml:"RecipeId"`
	Count     float64                  `yaml:"Count"`

	// To be loaded at runtime from the recipe lister output data
	Machine *recipe_lister.AssemblingMachine
	Recipe  *recipe_lister.Recipe
}
type ProcessList struct {
	Processes []ProcessDescription `yaml:"Processes"`
}

func (p *ProcessList) CalculateRates() []recipe_lister.RecipeRates {
	resp := make([]recipe_lister.RecipeRates, len(p.Processes))
	for i, proc := range p.Processes {
		resp[i] = proc.Recipe.CalculateRates(proc.Machine, 1, 1)
		for j, r := range resp[i].Inputs {
			resp[i].Inputs[j] = r * proc.Count
		}
		for j, r := range resp[i].Outputs {
			resp[i].Outputs[j] = r * proc.Count
		}
	}
	return resp
}

func LoadProcessList(file string) (*ProcessList, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("reading process list file: %w", err)
	}

	var list ProcessList
	err = yaml.Unmarshal(b, &list)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling process list: %w", err)
	}

	return &list, nil
}

func (p *ProcessList) PopulateData(data *recipe_lister.GameData) error {
	for i := range p.Processes {
		recipe, ok := data.Recipes[p.Processes[i].RecipeId]
		if !ok {
			return fmt.Errorf("recipe %s not found", p.Processes[i].RecipeId)
		}
		machine, ok := data.Machines[p.Processes[i].MachineId]
		if !ok {
			return fmt.Errorf("machine %s not found", p.Processes[i].MachineId)
		}
		p.Processes[i].Recipe = &recipe
		p.Processes[i].Machine = &machine
	}
	return nil
}
