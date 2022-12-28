package main

import (
	"time"

	"github.com/yzaimoglu/election/updater/controllers"
)

func main() {
	loop := true
	for loop {
		controllers.LoopQuarter()
		controllers.LoopDistrict()
		controllers.LoopConstituency()
		controllers.LoopCity()
		time.Sleep(10 * time.Minute)
	}
}
