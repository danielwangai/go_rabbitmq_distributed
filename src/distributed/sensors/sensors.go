package main

import (
	"flag"
)

var name = flag.String("name", "sensor", "Name of the sensor")
var freq = flag.Uint("freq", 5, "Update frequency in cycles/sec")
var max = flag.Float64("max", 5., "Maximum value for generated readings")
var min = flag.Float64("min", 5., "Minimum value for generated readings")
var stepSize = flag.Float64("step", 0.1, "Maximum allowable change per measurement.")

func main() {
	flag.Parse()
}
