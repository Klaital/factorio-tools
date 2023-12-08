package recipe_lister

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Generator struct {
	Name                string   `json:"name"`
	LocalisedName       []string `json:"localised_name"`
	MaximumTemperature  int      `json:"maximum_temperature"`
	Effectivity         int      `json:"effectivity"`
	FluidUsagePerTick   float64  `json:"fluid_usage_per_tick"`
	MaxEnergyProduction int      `json:"max_energy_production"`
	FriendlyMapColor    struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
		A int `json:"a"`
	} `json:"friendly_map_color"`
	EnemyMapColor struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
		A int `json:"a"`
	} `json:"enemy_map_color"`
	EnergySource struct {
		Electric struct {
			Drain     int     `json:"drain"`
			Emissions float64 `json:"emissions"`
		} `json:"electric"`
	} `json:"energy_source"`
	Pollution float64 `json:"pollution"`
}

type Boiler struct {
	Name              string   `json:"name"`
	LocalisedName     []string `json:"localised_name"`
	MaxEnergyUsage    int      `json:"max_energy_usage"`
	TargetTemperature int      `json:"target_temperature"`
	FriendlyMapColor  struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
		A int `json:"a"`
	} `json:"friendly_map_color"`
	EnemyMapColor struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
		A int `json:"a"`
	} `json:"enemy_map_color"`
	EnergySource struct {
		Electric struct {
			Emissions              float64 `json:"emissions"`
			MaxTemperature         int     `json:"max_temperature"`
			DefaultTemperature     int     `json:"default_temperature"`
			SpecificHeat           int     `json:"specific_heat"`
			MaxTransfer            float64 `json:"max_transfer"`
			MinTemperatureGradient int     `json:"min_temperature_gradient"`
			MinWorkingTemperature  int     `json:"min_working_temperature"`
		} `json:"electric"`
	} `json:"energy_source"`
	Pollution float64 `json:"pollution"`
}

func LoadGenerators(directory string) (map[string]Generator, error) {
	f, err := os.Open(fmt.Sprintf("%s/generator.json", directory))
	if err != nil {
		return nil, fmt.Errorf("opening generators file: %w", err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading generators file: %w", err)
	}
	generators := make(map[string]Generator, 0)
	err = json.Unmarshal(b, &generators)
	if err != nil {
		return nil, fmt.Errorf("parsing generators file: %w", err)
	}

	// Success!
	return generators, nil
}
func LoadBoilers(directory string) (map[string]Boiler, error) {
	f, err := os.Open(fmt.Sprintf("%s/boiler.json", directory))
	if err != nil {
		return nil, fmt.Errorf("opening boilers file: %w", err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading boilers file: %w", err)
	}
	data := make(map[string]Boiler, 0)
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, fmt.Errorf("parsing boilers file: %w", err)
	}

	// Success!
	return data, nil
}

type Reactor struct {
	Name             string   `json:"name"`
	LocalisedName    []string `json:"localised_name"`
	MaxEnergyUsage   int      `json:"max_energy_usage"`
	NeighbourBonus   float64  `json:"neighbour_bonus"`
	FriendlyMapColor struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
		A int `json:"a"`
	} `json:"friendly_map_color"`
	EnemyMapColor struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
		A int `json:"a"`
	} `json:"enemy_map_color"`
	EnergySource struct {
		Fluid struct {
			Emissions          float64 `json:"emissions"`
			Effectivity        int     `json:"effectivity"`
			BurnsFluid         bool    `json:"burns_fluid"`
			ScaleFluidUsage    bool    `json:"scale_fluid_usage"`
			FluidUsagePerTick  int     `json:"fluid_usage_per_tick"`
			MaximumTemperature int     `json:"maximum_temperature"`
			FluidBox           struct {
				Index          int    `json:"index"`
				ProductionType string `json:"production_type"`
				BaseArea       int    `json:"base_area"`
				BaseLevel      int    `json:"base_level"`
				Height         int    `json:"height"`
				Volume         int    `json:"volume"`
			} `json:"fluid_box"`
		} `json:"fluid"`
	} `json:"energy_source"`
	Pollution float64 `json:"pollution"`
}

func LoadReactors(directory string) (map[string]Reactor, error) {
	f, err := os.Open(fmt.Sprintf("%s/reactor.json", directory))
	if err != nil {
		return nil, fmt.Errorf("opening reactors file: %w", err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading reactors file: %w", err)
	}
	data := make(map[string]Reactor, 0)
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, fmt.Errorf("parsing reactors file: %w", err)
	}

	// Success!
	return data, nil
}

// PowerOut calculates the actual power yield in watts of heat, factoring in the neighbor bonus
func (r Reactor) PowerOut(neighbors int) int {
	multiplier := 1.0 + (float64(neighbors) * r.NeighbourBonus)
	power := multiplier * float64(r.MaxEnergyUsage)
	return int(power)
}
