package main

type Calculator func(cfg CalculationConfig, reactorCount int) CalculationResults

func Calc3xN(cfg CalculationConfig, reactorCount int) CalculationResults {
	res := CalculationResults{}
	// round down the number of reactors to the nearest multiple of 3
	rowsOfReactors := reactorCount / 3
	extraReactors := reactorCount % 3
	actualReactorCount := reactorCount - extraReactors
	res.ReactorQty = actualReactorCount

	// calculate the amount of fuel that will be consumed by the reactors
	res.FuelConsumed = (int64(cfg.Reactor.MaxEnergyUsage) * int64(actualReactorCount)) / 1000000

	// calculate the total heat generated, including neighbor bonuses
	cornerPower := 4 * cfg.Reactor.PowerOut(2)
	edgePower := (rowsOfReactors * 2) * cfg.Reactor.PowerOut(3)
	interiorPower := rowsOfReactors * cfg.Reactor.PowerOut(4)
	totalPower := int64(cornerPower) + int64(edgePower) + int64(interiorPower)
	res.MwYield = totalPower / 1000000

	// calculate the number of heat exchangers needed
	hex := totalPower / int64(cfg.Boiler.MaxEnergyUsage)
	res.BoilerQty = int(hex)

	// calculate the amount of steam produced
	// water consumption doesn't seem to be included the in the data dump.
	// I'm hardcoding it here after looking up the basic heat exchanger in-game.
	res.WaterRequired = 120 * res.BoilerQty

	// calculate the number of turbines needed
	turbines := totalPower / int64(cfg.Generator.MaxEnergyProduction)
	res.GeneratorQty = int(turbines)

	// TODO: auto-select the turbine based on the selected heat exchanger temperature.

	return res
}
