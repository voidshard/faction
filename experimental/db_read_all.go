package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/config"

	"github.com/voidshard/faction/internal/db"
)

func main() {
	// Reads everything written by db_insert_all.go
	// Nb. we don't iter with any token here, which generally would be advised

	cfg := &config.Database{
		Driver:   config.DatabaseSQLite3,
		Name:     "test.sqlite",
		Location: "/tmp",
	}
	conn, err := db.New(cfg)
	if err != nil {
		panic(err)
	}

	perr := func(err error) {
		if err != nil {
			fmt.Println(err)
		}
	}

	tick, _ := conn.Tick()
	fmt.Printf("tick %d\n", tick)

	fmt.Println("areas")
	areas, _, err := conn.Areas("", nil)
	perr(err)
	for _, i := range areas {
		fmt.Printf("\t%v\n", i)
	}

	fmt.Println("routes")
	routes, _, err := conn.Routes("", nil)
	perr(err)
	for _, i := range routes {
		fmt.Printf("\t%v\n", i)
	}

	fmt.Println("governments")
	govts, _, err := conn.Governments("", nil)
	perr(err)
	for _, i := range govts {
		fmt.Printf("\t%v\n\toutlawed\n", i)
		for item, illegal := range i.Outlawed.Commodities {
			fmt.Printf("\t\t%v: %v\n", item, illegal)
		}
		for action, illegal := range i.Outlawed.Actions {
			fmt.Printf("\t\t%v: %v\n", action, illegal)
		}
		for faction, illegal := range i.Outlawed.Factions {
			fmt.Printf("\t\t%v: %v\n", faction, illegal)
		}
	}

	fmt.Println("factions")
	factions, _, err := conn.Factions("", nil)
	perr(err)
	for _, i := range factions {
		fmt.Printf("\t%v\n", i)
	}

	fmt.Println("people")
	people, _, err := conn.People("", nil)
	perr(err)
	for _, i := range people {
		fmt.Printf("\t%v\n", i)
	}

	fmt.Println("families")
	families, _, err := conn.Families("", nil)
	perr(err)
	for _, i := range families {
		fmt.Printf("\t%v\n", i)
	}

	fmt.Println("plots")
	plots, _, err := conn.Plots("", nil)
	perr(err)
	for _, i := range plots {
		fmt.Printf("\t%v\n", i)
	}

	fmt.Println("jobs")
	jobs, _, err := conn.Jobs("", nil)
	perr(err)
	for _, i := range jobs {
		fmt.Printf("\t%v\n", i)
	}

	for _, r := range []db.Relation{
		db.RelationPersonPersonRelationship,
		db.RelationPersonPersonTrust,
		db.RelationPersonFactionAffiliation,
	} {
		fmt.Println("tuple", r)
		rel, _, err := conn.Tuples(r, "", nil)
		perr(err)
		for _, i := range rel {
			fmt.Printf("\t%v\n", i)
		}
	}

	for _, r := range []db.Relation{db.RelationPersonPersonTrust} {
		fmt.Println("modifier", r)
		rel, _, err := conn.Modifiers(r, "", nil)
		perr(err)
		for _, i := range rel {
			fmt.Printf("\t%v\n", i)
		}
	}

	fmt.Println("modifier sums", db.RelationPersonPersonTrust)
	sums, _, err := conn.ModifiersSum(db.RelationPersonPersonTrust, "", nil)
	perr(err)
	for _, i := range sums {
		fmt.Printf("\t%v\n", i)
	}

}
