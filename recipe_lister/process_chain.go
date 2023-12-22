package recipe_lister

import (
	"fmt"
	"log/slog"
)

type Process struct {
	ID            string
	Recipe        Recipe
	Machine       AssemblingMachine
	MachineCount  float64
	Modules       ModuleConfig
	BeaconModules ModuleConfig
	Parent        struct {
		ID          string
		ComponentID ItemName
	}
}
type ProcessChain struct {
	OutputTargetRates map[string]float64 // how much per second to produce
	Processes         []Process
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
