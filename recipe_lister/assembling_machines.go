package recipe_lister

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Machine interface {
	GetOperatingWatts() float64
	GetOperatingKiloWatts() float64
	GetIdleWatts() float64
}

type AssemblingMachine struct {
	Name string `json:"name"`
	EnergyUsage float64 `json:"energy_usage"`
	Drain float64 `json:"drain"`
	CraftingSpeed float64 `json:"crafting_speed"`
	ModuleInventorySize int64 `json:"module_inventory_size"`
}

func (m AssemblingMachine) GetOperatingWatts() float64 {
	return m.EnergyUsage
}
func (m AssemblingMachine) GetOperatingKiloWatts() float64 {
	return m.EnergyUsage / 1000.0
}
func (m AssemblingMachine) GetIdleWatts() float64 {
	return m.Drain
}

type Inserter struct {
	Name string `json:"name"`
	EnergyUsage float64 `json:"max_energy_usage"`
	Drain float64 `json:"drain"`
}
func (m Inserter) GetOperatingWatts() float64 {
	return m.EnergyUsage
}
func (m Inserter) GetOperatingKiloWatts() float64 {
	return m.EnergyUsage / 1000.0
}
func (m Inserter) GetIdleWatts() float64 {
	return m.Drain
}

func LoadAssemblingMachinesFile(path string) (dataSet map[string]AssemblingMachine, err error) {
	assemblingMachines := make(map[string]AssemblingMachine, 0)
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileBytes, &assemblingMachines)
	if err != nil {
		return nil, err
	}

	return assemblingMachines, nil
}

func LoadMachinesDirectory(path string) (dataSet map[string]Machine, err error) {
	assemblingMachines, err := LoadAssemblingMachinesFile(fmt.Sprintf("%s/assembling-machine.json", path))
	if err != nil {
		return nil, err
	}
	inserters := make(map[string]Inserter, 0)
	fileBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/inserter.json", path))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileBytes, &inserters)
	if err != nil {
		return nil, err
	}

	// Take the union of all machines
	dataSet = make(map[string]Machine)
	for _, machine := range assemblingMachines {
		dataSet[machine.Name] = machine
	}
	for _, machine := range inserters {
		dataSet[machine.Name] = machine
	}

	return dataSet, nil
}
