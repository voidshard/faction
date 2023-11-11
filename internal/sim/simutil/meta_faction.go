package simutil

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/structs"
)

// metaFaction is a working data set for operations on factions + associated data
type MetaFaction struct {
	Faction         *structs.Faction
	Actions         []string
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
		Actions:         []string{},
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
	err := conn.SetFactions(f.Faction)
	if err != nil {
		return err
	}
	err = conn.SetPlots(unique(f.Plots)...)
	if err != nil {
		return err
	}
	err = conn.SetEvents(unique(f.Events)...)
	if err != nil {
		return err
	}
	err = conn.SetTuples(db.RelationFactionCommodityImport, unique(f.Imports)...)
	if err != nil {
		return err
	}
	err = conn.SetTuples(db.RelationFactionCommodityExport, unique(f.Exports)...)
	if err != nil {
		return err
	}
	err = conn.SetTuples(db.RelationFactionTopicResearchWeight, unique(f.ResearchWeights)...)
	if err != nil {
		return err
	}
	err = conn.SetTuples(db.RelationFactionActionTypeWeight, unique(f.ActionWeights)...)
	if err != nil {
		return err
	}
	return conn.SetTuples(db.RelationFactionProfessionWeight, unique(f.ProfWeights)...)
}
