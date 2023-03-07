package db

import (
	"fmt"

	"github.com/voidshard/faction/pkg/structs"
)

// FactionDB is a helper struct that extends a Database implementation
// with helpful additional functions.
//
// These functions are useful to callers, but can be supplied regardless of
// the implementation assuming that the Database interface is met.
//
// Ie. all implementations of Database would supply these helper functions
// by doing the same thing, so we don't have to ask implementors to do it,
// rather we simply embed the Database into ourselves and build atop it.
type FactionDB struct {
	Database
}

// InTransaction is a helper function that adds automatic commit / rollback
// depending on if an error is returned.
func (f *FactionDB) InTransaction(do func(tx ReaderWriter) error) error {
	tx, err := f.Database.Transaction()
	if err != nil {
		return nil
	}

	err = do(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ChangeGoverningFaction is a helper function that changes the governing faction of some area(s)
// and any LandRight(s) they contain to the given governing faction.
func (f *FactionDB) ChangeGoverningFaction(govID string, areaIDs []string) error {
	if !structs.IsValidID(govID) {
		return fmt.Errorf("invalid faction id: %s", govID)
	}
	for _, a := range areaIDs {
		if !structs.IsValidID(a) {
			return fmt.Errorf("invalid area id: %s", a)
		}

		// TODO: we could make this more efficient.
		//
		// We do it this way because although we know there will be
		// at most one area per ID, there could be any number of land rights.
		//
		// In this fashion we make sure a given area is consistent (an area
		// and any rights it contains will have the same governing faction),
		// but we don't apply the entire change government operation in a single
		// transaction.

		areas, _, err := f.Areas("", &AreaFilter{ID: a})
		if err != nil {
			return err
		}
		if len(areas) != 1 {
			return fmt.Errorf("area %s not found", a)
		}

		areas[0].GoverningFactionID = govID
		lf := &LandRightFilter{AreaID: a}

		err = f.InTransaction(func(tx ReaderWriter) error {
			err = tx.SetAreas(areas[0])
			if err != nil {
				return nil
			}

			var (
				land  []*structs.LandRight
				token string
			)
			for {
				land, token, err = tx.LandRights(token, lf)
				if err != nil {
					return err
				}
				for _, l := range land {
					l.GoverningFactionID = govID
				}
				err = tx.SetLandRights(land...)
				if err != nil {
					return err
				}
				if token == "" {
					break
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// NewPump returns a pump - a buffered writer for the database.
func (f *FactionDB) NewPump() *Pump {
	return newPump(f)
}
