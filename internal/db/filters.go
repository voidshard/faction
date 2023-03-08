package db

import (
	"github.com/voidshard/faction/pkg/structs"
)

type AreaFilter struct {
	ID                 string
	GoverningFactionID string
}

type GovernmentFilter struct {
	ID string
}

type RouteFilter struct {
	SourceAreaID string
	TargetAreaID string
}

type EthosFilter struct {
	// provide some bounds for ethos matches
	MinEthos *structs.Ethos
	MaxEthos *structs.Ethos
}

type PersonFilter struct {
	EthosFilter

	ID            string
	JobID         string
	AreaID        string
	BirthFamilyID string

	Race string

	IncludeDead bool // by default we do not infact see dead people
}

type PlotFilter struct {
	ID             string
	OwnerFactionID string
	AreaID         string
}

type FactionFilter struct {
	EthosFilter

	ID              string
	GovernmentID    string
	ReligionID      string
	ParentFactionID string
}

type FamilyFilter struct {
	ID               string
	BirthFactionID   string
	FactionID        string
	OnlyChildBearing bool // restricts the search to only child bearing families
	MaleID           string
	FemaleID         string
}

type JobFilter struct {
	ID              string
	SourceFactionID string
	SourceAreaID    string
	TargetAreaID    string
	TargetMetaKey   structs.MetaKey
	TargetMetaVal   string
	MinSecrecy      int
	MaxSecrecy      int // values <= 0 are ignored
	State           structs.JobState
	TickEndsBefore  int
}

type LandRightFilter struct {
	ID     string
	AreaID string

	GoverningFactionID   string
	ControllingFactionID string
}

type TupleFilter struct {
	Subject string
	Object  string
}

type ModifierFilter struct {
	TupleFilter

	MinTickExpires int // values <= 0 are ignored
	MaxTickExpires int // values <= 0 are ignored

	MetaKey structs.MetaKey
	MetaVal string
}
