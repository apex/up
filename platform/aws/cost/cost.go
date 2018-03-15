// Package cost provides utilities for calculating AWS Lambda pricing.
package cost

// pricePerInvoke is the cost per function invocation.
var pricePerInvoke = 0.0000002

// pricePerRequestUnit is the cost per api gateway request unit.
var pricePerRequestUnit = 5

// requestUnit is 5 million requests.
var requestUnit = 5e6

// memoryConfigurations available.
var memoryConfigurations = map[int]float64{
	128:  0.000000208,
	192:  0.000000313,
	256:  0.000000417,
	320:  0.000000521,
	384:  0.000000625,
	448:  0.000000729,
	512:  0.000000834,
	576:  0.000000938,
	640:  0.000001042,
	704:  0.000001146,
	768:  0.00000125,
	832:  0.000001354,
	896:  0.000001459,
	960:  0.000001563,
	1024: 0.000001667,
	1088: 0.000001771,
	1152: 0.000001875,
	1216: 0.00000198,
	1280: 0.000002084,
	1344: 0.000002188,
	1408: 0.000002292,
	1472: 0.000002396,
	1536: 0.000002501,
}

// Requests returns the cost for the given number of http requests.
func Requests(n int) float64 {
	return (float64(n) / float64(requestUnit)) * float64(pricePerRequestUnit)
}

// Rate returns the cost per 100ms for the given `memory` configuration in megabytes.
func Rate(memory int) float64 {
	return memoryConfigurations[memory]
}

// Invocations returns the cost of `n` requests.
func Invocations(n int) float64 {
	return pricePerInvoke * float64(n)
}

// Duration returns the cost of `ms` for the given `memory` configuration in megabytes.
func Duration(ms, memory int) float64 {
	return Rate(memory) * (float64(ms) / 100)
}
