package base

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/structs"
)

// metaFaction is a working data set for operations on factions + associated data
type metaFaction struct {
	faction         *structs.Faction
	actions         []structs.ActionType
	actionWeights   []*structs.Tuple
	plots           []*structs.Plot
	profWeights     []*structs.Tuple
	researchWeights []*structs.Tuple
	areas           map[string]bool
}

func writeMetaFaction(conn *db.FactionDB, f *metaFaction) error {
	return conn.InTransaction(func(tx db.ReaderWriter) error {
		err := tx.SetFactions(f.faction)
		if err != nil {
			return err
		}
		err = tx.SetPlots(f.plots...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationFactionTopicResearchWeight, f.researchWeights...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationFactionActionTypeWeight, f.actionWeights...)
		if err != nil {
			return err
		}
		return tx.SetTuples(db.RelationFactionProfessionWeight, f.profWeights...)
	})
}
