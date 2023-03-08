package main

import (
	"fmt"

	"math"

	"github.com/voidshard/faction/pkg/structs"
)

func main() {

	e0 := (&structs.Ethos{}).Sub(math.MaxInt)
	e1 := (&structs.Ethos{}).Add(math.MaxInt)

	fmt.Println(structs.EthosDistance(e0, e1))
}
