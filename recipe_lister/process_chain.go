package recipe_lister

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"math"
	"os"
)

type ParentConfig struct {
	ID          string   `yaml:"id"`
	ComponentID ItemName `yaml:"component"`
}
type Process struct {
	ID            string            `yaml:"id"`
	Recipe        Recipe            `yaml:"recipe"`
	Machine       AssemblingMachine `yaml:"machine"`
	MachineCount  float64           `yaml:"machinecount"`
	Modules       ModuleConfig      `yaml:"modules"`
	BeaconModules ModuleConfig      `yaml:"beaconmodules"`
	Parent        ParentConfig      `yaml:"parent"`
}
type ProcessChain struct {
	OutputTargetRates map[string]float64 `yaml:"OutputTargetRates"` // how much per second to produce
	Processes         []Process          `yaml:"Processes"`
}

func LoadProcessChain(processFile string, recipeListerDir string) (*ProcessChain, error) {
	// Load the process list itself
	b, err := os.ReadFile(processFile)
	if err != nil {
		return nil, fmt.Errorf("loading process file: %w", err)
	}

	var processes ProcessChain
	err = yaml.Unmarshal(b, &processes)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling process file: %w", err)
	}

	// Load the recipe and machine data
	recipes, err := LoadRecipes(recipeListerDir)
	if err != nil {
		return nil, err
	}
	machines, err := LoadAllBuilders(recipeListerDir)
	if err != nil {
		return nil, err
	}

	// Populate the process chain with game data
	processes.AnnotateGameData(recipes, machines)
	return &processes, nil
}

func (c *ProcessChain) AnnotateGameData(recipes map[RecipeName]Recipe, machines map[MachineName]AssemblingMachine) {
	for i, process := range c.Processes {
		c.Processes[i].Recipe = recipes[process.Recipe.Name]
		c.Processes[i].Machine = machines[process.Machine.Name]
	}
}

func (c *ProcessChain) GetProcessById(id string) *Process {
	for i, p := range c.Processes {
		if p.ID == id {
			return &c.Processes[i]
		}
	}
	return nil
}

// ComputeMachineCounts iterates over the processes which have a parent,
// and calculates how many are needed in order to satisfy the parent's
// input requirements for the specific item.
func (c *ProcessChain) ComputeMachineCounts() error {
	for i, p := range c.Processes {
		if len(p.Parent.ID) == 0 {
			continue
		}
		if len(p.Parent.ComponentID) == 0 {
			return fmt.Errorf("computing counts for process %s: no component ID given for parent %s", p.ID, p.Parent.ID)
		}

		parentProcess := c.GetProcessById(p.Parent.ID)
		c.Processes[i].MatchProduction(parentProcess, p.Parent.ComponentID)
	}

	// Success!
	return nil
}

func (c *ProcessChain) FindProcessByID(id string) *Process {
	for i, p := range c.Processes {
		if p.ID == id {
			return &c.Processes[i]
		}
	}
	// Not found
	return nil
}

func (p *Process) SecondsPerCycle() float64 {
	return p.Recipe.Energy / p.Machine.CraftingSpeed
}
func (p *Process) ItemsPerCyclePerMachine() RecipeRates {
	resp := NewRates()
	// TODO: account for productivity modules
	for _, item := range p.Recipe.Ingredients {
		probability := item.Probability
		if probability == 0 {
			probability = 1
		}
		resp.Inputs[item.Name] = item.Amount * probability
	}
	for _, item := range p.Recipe.Products {
		probability := item.Probability
		if probability == 0 {
			probability = 1
		}
		amountMin := item.AmountMin
		amountMax := item.AmountMax
		if item.Amount > 0 {
			amountMin = item.Amount
			amountMax = item.Amount
		}
		qty := (amountMin + amountMax) / 2.0
		resp.Outputs[item.Name] = qty * probability
	}
	return resp
}

// ItemsPerSecondPerMachine converts the Items per cycle into per-second
func (p *Process) ItemsPerSecondPerMachine() RecipeRates {
	cycleRates := p.ItemsPerCyclePerMachine()
	cycleRates.ModifyAll(func(x float64) float64 {
		return x / p.SecondsPerCycle()
	})
	return cycleRates
}

// ItemsPerSecond returns the final rates of production and consumption,
// factoring in modules, beacons, and the number of machines running
// the process
func (p *Process) ItemsPerSecond() RecipeRates {
	rates := p.ItemsPerSecondPerMachine()
	rates.ModifyAll(func(x float64) float64 {
		return x * p.MachineCount
	})
	return rates
}

// MatchProduction updates a process's number of machines in order to
// produce enough of the given item to satisfy the input requirements
// of the other process.
func (p *Process) MatchProduction(otherProcess *Process, name ItemName) {
	targetRates := otherProcess.ItemsPerSecond()
	targetRate, ok := targetRates.Inputs[name]
	if !ok {
		slog.Error("unable to match production rate. Parent does not consume item.", "item", name, "child_process", p.ID, "parent_process", otherProcess.ID)
		return
	}
	productionRates := p.ItemsPerSecondPerMachine()
	productionRate, ok := productionRates.Outputs[name]
	if !ok {
		slog.Error("unable to match production rate. Child does not produce item.", "item", name, "child_process", p.ID, "parent_process", otherProcess.ID)
		return
	}

	p.MachineCount = targetRate / productionRate
}

func (c *ProcessChain) TotalIO() RecipeRates {
	sum := NewRates()
	for _, process := range c.Processes {
		processRates := process.ItemsPerSecond()
		sum.Add(processRates)
	}
	return Split(sum.Merge())
}

// Merge combines input and outputs in one set, with inputs as negatives.
func (r *RecipeRates) Merge() map[ItemName]float64 {
	all := make(map[ItemName]float64, 0)
	for n, rate := range r.Inputs {
		all[n] = -1.0 * rate
	}
	for n, rate := range r.Outputs {
		all[n] = rate
	}
	return all
}

var epsilon float64 = 1e-6

// Split is the opposite of RecipeRates#Merge. Negative rates translate
// to inputs, positive to outputs. Rates less than epsilon are omitted.
func Split(rates map[ItemName]float64) RecipeRates {
	split := NewRates()
	for name, rate := range rates {
		abs := math.Abs(rate)
		if abs < epsilon {
			continue
		}
		if rate < 0 {
			split.Inputs[name] = abs
		} else {
			split.Outputs[name] = abs
		}
	}
	return split
}

func (r *RecipeRates) Add(more RecipeRates) {
	for name, rate := range more.Inputs {
		r.Inputs[name] = r.Inputs[name] + rate
	}
	for name, rate := range more.Outputs {
		r.Outputs[name] = r.Outputs[name] + rate
	}
}
