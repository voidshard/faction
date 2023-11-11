package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"

	"github.com/voidshard/faction/internal/db"
)

func main() {
	// Tests out the incr logic for tuples / modifiers
	conn, err := db.New(config.DefaultDatabase())
	if err != nil {
		panic(err)
	}

	subject := structs.NewID()
	object := structs.NewID()

	aff := &structs.Tuple{Subject: subject, Object: object, Value: 10}

	tx, err := conn.Transaction()
	if err != nil {
		panic(err)
	}
	errRollback := func(err error) {
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	err = tx.SetTuples(db.RelationPersonFactionAffiliation, aff)
	errRollback(err)
	fmt.Println("wrote tuple", db.RelationPersonFactionAffiliation, aff)

	fmt.Println(subject, object)
	q := db.Q(db.F(db.Subject, db.Equal, subject), db.F(db.Object, db.Equal, object))

	err = tx.IncrTuples(db.RelationPersonFactionAffiliation, 10, q)
	errRollback(err)
	fmt.Println("incremented tuple", db.RelationPersonFactionAffiliation, aff, "by 10")

	err = tx.IncrTuples(db.RelationPersonFactionAffiliation, -100, q)
	errRollback(err)
	fmt.Println("incremented tuple", db.RelationPersonFactionAffiliation, aff, "by -100")

	err = tx.Commit()
	errRollback(err)
}
