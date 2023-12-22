package recipe_lister

import (
	"embed"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"
)

//go:embed testdata/*.json
var fixtures embed.FS

func fixtureRecipes() map[RecipeName]Recipe {
	recipes := make(map[RecipeName]Recipe, 0)
	fixtureData, err := fixtures.ReadDir("testdata")
	if err != nil {
		panic("failed to read testdata fixtures")
	}
	for _, d := range fixtureData {
		if strings.HasSuffix(d.Name(), ".json") && strings.HasPrefix(d.Name(), "recipe_") {
			var tmp Recipe
			b, err := fixtures.ReadFile("testdata/" + d.Name())
			if err != nil {
				panic(fmt.Sprintf("failed to read fixture file %s: %s", d.Name(), err.Error()))
			}
			err = json.Unmarshal(b, &tmp)
			if err != nil {
				panic(fmt.Sprintf("failed to unmarshal fixture file %s: %s", d.Name(), err.Error()))
			}

			recipes[tmp.Name] = tmp
		}
	}
	return recipes
}
func fixtureMachines() map[MachineName]AssemblingMachine {
	machines := make(map[MachineName]AssemblingMachine, 0)
	fixtureData, err := fixtures.ReadDir("testdata")
	if err != nil {
		panic("failed to read testdata fixtures")
	}
	for _, d := range fixtureData {
		if strings.HasSuffix(d.Name(), ".json") && strings.HasPrefix(d.Name(), "machine_") {
			var tmp AssemblingMachine
			b, err := fixtures.ReadFile("testdata/" + d.Name())
			if err != nil {
				panic(fmt.Sprintf("failed to read fixture file %s: %s", d.Name(), err.Error()))
			}
			err = json.Unmarshal(b, &tmp)
			if err != nil {
				panic(fmt.Sprintf("failed to unmarshal fixture file %s: %s", d.Name(), err.Error()))
			}

			machines[tmp.Name] = tmp
		}
	}
	return machines
}

const float64EqualityThreshold = 1e-6

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func stringArrayIncludes(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
func equalRates(a, b RecipeRates) bool {
	// Check that the keys are the same
	if len(a.Inputs) != len(b.Inputs) || len(a.Outputs) != len(b.Outputs) {
		return false
	}
	keysA := make([]string, 0, len(a.Inputs))
	for k := range a.Inputs {
		keysA = append(keysA, string(k))
	}
	for k := range b.Inputs {
		if !stringArrayIncludes(keysA, string(k)) {
			return false
		}
	}
	keysA = make([]string, 0, len(a.Outputs))
	for k := range a.Outputs {
		keysA = append(keysA, string(k))
	}
	for k := range b.Outputs {
		if !stringArrayIncludes(keysA, string(k)) {
			return false
		}
	}

	// check that the values are all approximately the same
	for name, rate := range a.Inputs {
		if !almostEqual(rate, b.Inputs[name]) {
			return false
		}
	}
	for name, rate := range a.Outputs {
		if !almostEqual(rate, b.Outputs[name]) {
			return false
		}
	}

	// All the tests have passed
	return true
}
func TestProcess_SecondsPerCycle(t *testing.T) {
	allRecipes := fixtureRecipes()
	allMachines := fixtureMachines()
	cases := []struct {
		name           string
		process        Process
		expectedResult float64
	}{
		{
			name: "basic",
			process: Process{
				Recipe:  allRecipes["washing-1"],
				Machine: allMachines["washing-plant-2"],
			},
			expectedResult: 2.222222,
		},
		// TODO: add tests with modules and beacons
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.process.SecondsPerCycle()
			if !almostEqual(actual, tt.expectedResult) {
				t.Errorf("Incorrect cycle time. Expected %f, got %f", tt.expectedResult, actual)
			}
		})
	}
}

func TestProcess_ItemsPerCyclePerMachine(t *testing.T) {
	tests := []struct {
		name    string
		process Process
		want    RecipeRates
	}{
		{
			name: "washing-1",
			process: Process{
				Recipe:        fixtureRecipes()["washing-1"],
				Modules:       ModuleConfig{},
				BeaconModules: ModuleConfig{},
			},
			want: RecipeRates{
				Inputs: map[ItemName]float64{
					"water-viscous-mud": 200,
					"water":             50,
				},
				Outputs: map[ItemName]float64{
					"solid-mud":            0.75,
					"water-heavy-mud":      200,
					"gas-hydrogen-sulfide": 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.process.ItemsPerCyclePerMachine()
			if !equalRates(tt.want, actual) {
				t.Errorf("ItemsPerCyclePerMachine not correct. Expected %+v, got %+v", tt.want, actual)
			}
		})
	}
}

func TestProcess_ItemsPerSecondPerMachine(t *testing.T) {
	tests := []struct {
		name    string
		process Process
		want    RecipeRates
	}{
		{
			name: "washing-1",
			process: Process{
				Recipe:        fixtureRecipes()["washing-1"],
				Machine:       fixtureMachines()["washing-plant-2"],
				Modules:       ModuleConfig{},
				BeaconModules: ModuleConfig{},
				MachineCount:  1.0,
			},
			want: RecipeRates{
				Inputs: map[ItemName]float64{
					"water-viscous-mud": 90,
					"water":             22.5,
				},
				Outputs: map[ItemName]float64{
					"solid-mud":            0.3375,
					"water-heavy-mud":      90,
					"gas-hydrogen-sulfide": 0.899999,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.process.ItemsPerSecondPerMachine()
			if !equalRates(tt.want, actual) {
				t.Errorf("ItemsPerSecondPerMachine not correct. Expected %+v, got %+v", tt.want, actual)
			}
		})
	}
}

func TestProcess_MatchProduction(t *testing.T) {
	tests := []struct {
		name                 string
		child                Process
		parent               Process
		expectedMachineCount float64
	}{
		{
			name: "mud for soil",
			child: Process{
				Recipe:        fixtureRecipes()["washing-1"],
				Machine:       fixtureMachines()["washing-plant-2"],
				Modules:       ModuleConfig{},
				BeaconModules: ModuleConfig{},
			},
			parent: Process{
				Recipe:       fixtureRecipes()["solid-soil"],
				Machine:      fixtureMachines()["assembling-machine-2"],
				MachineCount: 1,
			},
			expectedMachineCount: 0.55555555,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.child.MatchProduction(&tt.parent, "solid-mud")
			if !almostEqual(tt.child.MachineCount, tt.expectedMachineCount) {
				t.Errorf("machine count not correct. Expected %f, got %f", tt.expectedMachineCount, tt.child.MachineCount)
			}
		})
	}
}
