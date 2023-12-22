package recipe_lister

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestProcessChain_ComputeMachineCounts(t *testing.T) {
	allMachines := fixtureMachines()
	allRecipes := fixtureRecipes()

	tests := []struct {
		name                 string
		initial              ProcessChain
		updatedProcessCounts map[string]float64
		wantErr              bool
	}{
		{
			name: "two-process chain",
			initial: ProcessChain{
				Processes: []Process{
					{
						ID:           "parent",
						Recipe:       allRecipes["solid-soil"],
						MachineCount: 1,
						Machine:      allMachines["assembling-machine-2"],
					},
					{
						ID:      "child",
						Recipe:  allRecipes["washing-1"],
						Machine: allMachines["washing-plant-2"],
						Parent: ParentConfig{
							ID:          "parent",
							ComponentID: ItemName("solid-mud"),
						},
					},
				},
			},
			updatedProcessCounts: map[string]float64{
				"child": 0.555556,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.initial.ComputeMachineCounts()
			if err != nil {
				t.Errorf("ComputeMachineCounts error: %+v", err)
				return
			}
			for _, updatedProc := range tt.initial.Processes {
				expectedCount, ok := tt.updatedProcessCounts[updatedProc.ID]
				if !ok {
					// not testing this process
					continue
				}
				if !almostEqual(expectedCount, updatedProc.MachineCount) {
					t.Errorf("Incorrectly updated machine count for process '%s'. Expected %f, got %f", updatedProc.ID, expectedCount, updatedProc.MachineCount)
				}
			}
		})
	}
}

func TestLoadProcessChain(t *testing.T) {
	allMachines := fixtureMachines()
	allRecipes := fixtureRecipes()
	processes := ProcessChain{
		Processes: []Process{
			{
				Recipe:        allRecipes["washing-1"],
				Machine:       allMachines["washing-plant-2"],
				Modules:       ModuleConfig{},
				BeaconModules: ModuleConfig{},
				ID:            "mud production",
			}, {
				ID:           "soil manufacturing",
				Recipe:       allRecipes["solid-soil"],
				Machine:      allMachines["assembling-machine-2"],
				MachineCount: 1,
			},
		},
	}

	b, _ := yaml.Marshal(processes)
	fmt.Printf(string(b))
}
