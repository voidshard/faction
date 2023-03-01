package db

import (
	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/pkg/structs"

	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// sqlOperator is something that can perform an sql operation read/write
// We do this so we can have some lower level funcs that perform the query logic regardless
// of whether we are in a transaction or not.
type sqlOperator interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	Select(dest interface{}, query string, args ...interface{}) error
}

// mstruct is a row of metadata
type mstruct struct {
	ID  string `db:"id"`
	Str string `db:"str"`
	Int int    `db:"int"`
}

func areas(op sqlOperator, token string, in []*AreaFilter) ([]*structs.Area, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, err
	}

}

// setAreas saves area information to the database
func setAreas(op sqlOperator, areas []*structs.Area) error {
	if len(areas) == 0 {
		return nil
	}

	qstr := fmt.Sprintf(`INSERT INTO %s (id, name, description, parent_id, type)
        VALUES (:id, :governing_faction) 
        ON CONFLICT (id) DO UPDATE SET
            governing_faction=EXCLUDED.governing_faction
        ;`,
		tableAreas,
	)

	_, err := op.NamedExec(qstr, areas)
	return err
}

// meta returns some metadata, if set
func meta(op sqlOperator, id string) (string, int, error) {
	if !dbutils.IsValidName(id) {
		return "", 0, fmt.Errorf("metadata key %s is invalid", id)
	}

	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE id=$1 LIMIT 1;",
		tableMeta,
	)

	result := []*mstruct{}
	err := op.Select(&result, query, id)
	if err != nil || len(result) == 0 {
		return "", 0, err
	}

	return result[0].Str, result[0].Int, nil
}

// setMeta sets some data in our meta table
func setMeta(op sqlOperator, id, strv string, intv int) error {
	if !dbutils.IsValidName(id) {
		return fmt.Errorf("metadata key %s is invalid", id)
	}

	// update schema version to current
	qstr := fmt.Sprintf(`INSERT INTO %s (id, str, int)
		VALUES (:id, :str, :int) 
		ON CONFLICT (id) DO UPDATE SET
		    int=EXCLUDED.int,
		    str=EXCLUDED.str
		;`,
		tableMeta,
	)
	_, err := op.NamedExec(qstr, map[string]interface{}{
		"id":  id,
		"str": strv,
		"int": intv,
	})
	return err
}

// queryByIds returns a query string given IDs.
func queryByIds(ids []string) (string, []interface{}) {
	if ids == nil || len(ids) == 0 {
		return "", nil
	}
	args := []interface{}{}
	or := []string{}
	for i, id := range ids {
		if !dbutils.IsValidID(id) {
			continue
		}
		name := fmt.Sprintf("$%d", i)
		args = append(args, id)
		or = append(or, fmt.Sprintf("id=:%s", name))
	}
	if len(or) == 0 {
		return "", nil
	}
	return fmt.Sprintf(" WHERE %s", strings.Join(or, " OR ")), args
}
