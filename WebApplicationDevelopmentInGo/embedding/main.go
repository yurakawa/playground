package main

import "fmt"

type Dog struct{}

func (d Dog) Bark() string {
	return "Woof!"
}

type BullDog struct{ Dog }

type ShibaInu struct{ Dog }

func (s *ShibaInu) Bark() string {
	return "Kyan!"
}

func DogVoice(d *Dog) string {
	return d.Bark()
}

func main() {
	bd := &BullDog{}
	fmt.Println(bd.Bark())
	si := &ShibaInu{}
	fmt.Println(si.Bark())

	// cannot use si (type *ShibaInu) as type *Dog in argument to DogVoice
	// fmt.Println(DogVoice(&Dog))
}
