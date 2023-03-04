package db

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
