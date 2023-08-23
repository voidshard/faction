package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/premade/fantasy"
	"github.com/voidshard/faction/pkg/structs"
)

func main() {
	actions := fantasy.Actions()

	f := &structs.Ethos{
		Pacifism: structs.MaxEthos / 2,
		Piety:    structs.MaxEthos,
		Ambition: structs.MaxEthos / 5,
		//		Piety:    structs.MaxEthos,
		//Ambition: structs.MinEthos / 2,
	}

	act, _ := actions[structs.ActionTypeWar]

	fmt.Println(&act.Ethos, f)
	dist := structs.EthosDistance(f, &act.Ethos)
	fmt.Println(1 - dist)

}
