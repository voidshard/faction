package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"

	"github.com/voidshard/faction/internal/db"
)

func main() {
	cfg := &config.Database{
		Driver:   config.DatabaseSQLite3,
		Name:     "test.sqlite",
		Location: "/tmp",
	}
	conn, err := db.New(cfg)
	if err != nil {
		panic(err)
	}

	area1 := &structs.Area{ID: structs.NewID()}
	area2 := &structs.Area{ID: structs.NewID()}

	tx, err := conn.Transaction()
	if err != nil {
		panic(err)
	}

	err = tx.SetAreas(area1, area2)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	allAreas, token, err := conn.Areas("", &db.AreaFilter{ID: area1.ID})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Areas: %v %s\n", allAreas, token)
}
