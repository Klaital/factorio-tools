package recipe_lister

type ModuleType string

const (
	PRODUCTIVITY ModuleType = "prod"
	SPEED        ModuleType = "speed"
)

// speedMultipliers is just hardcoded because it's easier. This is from Seablock / BobAngels modpack.
var speedMultipliers = map[ModuleType]map[int]float64{
	SPEED: {
		0: 1.2,
		1: 1.3,
		2: 1.5,
		3: 1.7,
	},
	PRODUCTIVITY: {
		0: 0.9,
		1: 0.88,
		2: 0.86,
		3: 0.85,
	},
}

type ModuleConfig struct {
	Module     ModuleType
	Level      int
	Count      int
	FromBeacon bool
}

func (m ModuleConfig) SpeedMultiplier() float64 {
	return speedMultipliers[m.Module][m.Count]
}
