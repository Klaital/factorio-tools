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

type Builder interface {
	GetName() MachineName
	GetOperatingWatts() float64
	GetOperatingKiloWatts() float64
	GetIdleWatts() float64
	SupportsCraftingCategory(categoryName string) bool
	GetCraftingSpeed() float64
	GetModuleInventoryCount() int64
}

type MachineName string

type AssemblingMachine struct {
	Name                MachineName     `json:"name" yaml:"name"`
	EnergyUsage         float64         `json:"energy_usage"`
	Drain               float64         `json:"drain"`
	CraftingSpeed       float64         `json:"crafting_speed"`
	ModuleInventorySize int64           `json:"module_inventory_size"`
	CraftingCategories  map[string]bool `json:"crafting_categories"`
}

func (m AssemblingMachine) GetName() MachineName {
	return m.Name
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
func (m AssemblingMachine) SupportsCraftingCategory(categoryName string) bool {
	return m.CraftingCategories[categoryName]
}

func (m AssemblingMachine) GetCraftingSpeed() float64 {
	return m.CraftingSpeed
}
func (m AssemblingMachine) GetModuleInventoryCount() int64 {
	return m.ModuleInventorySize
}

type Inserter struct {
	Name        MachineName `json:"name"`
	EnergyUsage float64     `json:"max_energy_usage"`
	Drain       float64     `json:"drain"`
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

func LoadAllBuilders(directory string) (map[MachineName]AssemblingMachine, error) {
	assemblers, err := LoadAssemblingMachinesFile(fmt.Sprintf("%s/assembling-machine.json", directory))
	if err != nil {
		return nil, fmt.Errorf("loading assembling machines: %w", err)
	}
	furnaces, err := LoadFurnacesFile(fmt.Sprintf("%s/furnace.json", directory))
	if err != nil {
		return nil, fmt.Errorf("loading furnaces: %w", err)
	}

	// Merge the two datasets
	machines := make(map[MachineName]AssemblingMachine, 0)
	for name := range assemblers {
		machines[assemblers[name].Name] = assemblers[name]
	}
	for name := range furnaces {
		machines[furnaces[name].Name] = furnaces[name]
	}

	return machines, nil
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

func LoadFurnacesFile(path string) (dataSet map[string]AssemblingMachine, err error) {
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

// LoadBuildersFromDirectory loads data on all machines capable of construction:
// furnaces, assembling machines, centrifuges, etc
func LoadBuildersFromDirectory(path string) (dataSet map[MachineName]Builder, err error) {
	assemblingMachines, err := LoadAssemblingMachinesFile(fmt.Sprintf("%s/assembling-machine.json", path))
	if err != nil {
		return nil, err
	}
	furnaces, err := LoadFurnacesFile(fmt.Sprintf("%s/furnace.json", path))
	if err != nil {
		return nil, err
	}
	// Take the union of all machines
	dataSet = make(map[MachineName]Builder)
	for _, machine := range assemblingMachines {
		dataSet[machine.Name] = machine
	}
	for _, machine := range furnaces {
		dataSet[machine.Name] = machine
	}
	return dataSet, nil
}

// LoadMachinesDirectory is used to load all power-consuming machines, suitable for calculating power requirements.
func LoadMachinesDirectory(path string) (dataSet map[MachineName]Machine, err error) {
	assemblingMachines, err := LoadAssemblingMachinesFile(fmt.Sprintf("%s/assembling-machine.json", path))
	if err != nil {
		return nil, err
	}
	furnaces, err := LoadFurnacesFile(fmt.Sprintf("%s/furnace.json", path))
	if err != nil {
		return nil, err
	}
	inserters := make(map[MachineName]Inserter, 0)
	fileBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/inserter.json", path))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileBytes, &inserters)
	if err != nil {
		return nil, err
	}

	// Take the union of all machines
	dataSet = make(map[MachineName]Machine)
	for _, machine := range assemblingMachines {
		dataSet[machine.Name] = machine
	}
	for _, machine := range inserters {
		dataSet[machine.Name] = machine
	}
	for _, machine := range furnaces {
		dataSet[machine.Name] = machine
	}

	return dataSet, nil
}
