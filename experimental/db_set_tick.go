package main

import (
	"flag"
	"fmt"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/config"
)

func main() {
	var tick int
	flag.IntVar(&tick, "tick", 1, "tick to set")
	flag.Parse()

	if tick < 1 {
		tick = 1
	}

	conn, err := db.New(config.DefaultDatabase())
	if err != nil {
		panic(err)
	}

	err = conn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetTick(tick)
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("tick set to", tick)
}
