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
	Imports         []*structs.Tuple
	Exports         []*structs.Tuple
	Areas           map[string]bool
	Events          []*structs.Event
}

func NewMetaFaction() *MetaFaction {
	return &MetaFaction{
		Faction:         &structs.Faction{},
		Actions:         []structs.ActionType{},
		ActionWeights:   []*structs.Tuple{},
		Plots:           []*structs.Plot{},
		ProfWeights:     []*structs.Tuple{},
		ResearchWeights: []*structs.Tuple{},
		Imports:         []*structs.Tuple{},
		Exports:         []*structs.Tuple{},
		Areas:           map[string]bool{},
		Events:          []*structs.Event{},
	}
}

func WriteMetaFaction(conn *db.FactionDB, f *MetaFaction) error {
	// de-dupe stuff before we enter a transaction.
	// By definition anything we "Set" overwrites any existing data so we only need to apply
	// the last object in the slice for each ID.
	// Technically the DB could resolve this (depending on the engine) but this is busy work
	// we can do out-of-transaction
	f.ActionWeights = unique(f.ActionWeights)
	f.Plots = unique(f.Plots)
	f.ProfWeights = unique(f.ProfWeights)
	f.ResearchWeights = unique(f.ResearchWeights)
	f.Imports = unique(f.Imports)
	f.Exports = unique(f.Exports)
	f.Events = unique(f.Events)
	return conn.InTransaction(func(tx db.ReaderWriter) error {
		err := tx.SetFactions(f.Faction)
		if err != nil {
			return err
		}
		err = tx.SetPlots(f.Plots...)
		if err != nil {
			return err
		}
		err = tx.SetEvents(f.Events...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationFactionCommodityImport, f.Imports...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationFactionCommodityExport, f.Exports...)
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
