package recipe_lister

import "testing"

func TestReactor_PowerOut(t *testing.T) {
	r := Reactor{
		Name:           "",
		LocalisedName:  nil,
		MaxEnergyUsage: 1000,
		NeighbourBonus: 0.125,
	}

	// expected neighbor bonus = 1 * (neighborcount * r.NeighborBonus)
	if 1125 != r.PowerOut(1) {
		t.Errorf("Incorrect power with 1 neighbor. Expected %d, got %d", 1125, r.PowerOut(1))
	}
}
