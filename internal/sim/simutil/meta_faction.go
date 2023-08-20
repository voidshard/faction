package simutil

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/structs"
)

// metaFaction is a working data set for operations on factions + associated data
type MetaFaction struct {
	Faction         *structs.Faction
	Actions         []structs.ActionType
	ActionWeights   []*structs.Tuple
	Plots           []*structs.Plot
	ProfWeights     []*structs.Tuple
	ResearchWeights []*structs.Tuple
	Areas           map[string]bool
}

func NewMetaFaction() *MetaFaction {
	return &MetaFaction{
		Faction:         &structs.Faction{},
		Actions:         []structs.ActionType{},
		ActionWeights:   []*structs.Tuple{},
		Plots:           []*structs.Plot{},
		ProfWeights:     []*structs.Tuple{},
		ResearchWeights: []*structs.Tuple{},
		Areas:           map[string]bool{},
	}
}

func WriteMetaFaction(conn *db.FactionDB, f *MetaFaction) error {
	return conn.InTransaction(func(tx db.ReaderWriter) error {
		err := tx.SetFactions(f.Faction)
		if err != nil {
			return err
		}
		err = tx.SetPlots(f.Plots...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationFactionTopicResearchWeight, f.ResearchWeights...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationFactionActionTypeWeight, f.ActionWeights...)
		if err != nil {
			return err
		}
		return tx.SetTuples(db.RelationFactionProfessionWeight, f.ProfWeights...)
	})
}
