package main

import (
	. "builder_pattern/car"
)


func main () {
	car := NewBuilder().Color(BlueColor).Wheels(SteelWheels).TopSpeed(KPH).Build()
	car.Drive()
	car.Stop()
}