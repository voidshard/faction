package db

import (
	"github.com/voidshard/faction/pkg/structs"
)

type AreaFilter struct {
	ID                 string
	GoverningFactionID string
}

type ethosFilter struct {
	// provide some bounds for ethos matches
	MinEthos *structs.Ethos
	MaxEthos *structs.Ethos
}

type PersonFilter struct {
	ethosFilter

	ID     string
	JobID  string
	AreaID string

	IncludeDead bool // by default we do not infact see dead people
}

type PlotFilter struct {
	ID             string
	OwnerFactionID string
	AreaID         string
}

type FactionFilter struct {
	ethosFilter

	ID              string
	GovernmentID    string
	ReligionID      bool
	ParentFactionID string
}

type FamilyFilter struct {
	ID             string
	IsChildBearing bool
	FactionID      string
}

type JobFilter struct {
	ID              string
	SourceFactionID string
	SourceAreaID    string
	TargetFactionID string
	TargetAreaID    string
	MinSecrecy      int
	MaxSecrecy      int
}

type LandRightFilter struct {
	ID     string
	AreaID string

	GoverningFactionID  string
	ControlledFactionID string
}

type TupleFilter struct {
	Subject string
	Object  string
}
