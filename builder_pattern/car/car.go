package car

import "fmt"

type carBuilder struct {
	color Color
	wheels Wheels
	speed Speed
}

type car struct {
	params carBuilder
}

func NewBuilder() *carBuilder {
	return &carBuilder{
		color: BlueColor,
		wheels: SportsWheels,
		speed: MPH,
	}
}

func (b *carBuilder) Color (color Color) Builder {
	b.color = color
	return b
}

func (b *carBuilder) Wheels (wheel Wheels) Builder {
	b.wheels = wheel
	return b
}

func (b *carBuilder) TopSpeed (speed Speed) Builder {
	b.speed = speed
	return b
}

func (b *carBuilder) Build () Interface {
	return &car {
		params: *b,
	}
}

func (c *car) Drive() error {
	fmt.Printf("Driveing: %#+v\n", c.params)
	return nil
}

func(c *car) Stop () error {
	fmt.Printf("Stop: %#+v\n", c.params)
	return nil
}
