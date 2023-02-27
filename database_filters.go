package faction

const (
	defaultLimit = 5000
)

type TupleFilter struct {
	Subject  *string
	Object   *string
	ValueMin *int
	ValueMax *int
}

type AreaFilter struct{}

type PersonFilter struct{}

type FactionFilter struct{}

type JobFilter struct{}

type LandRightFilter struct{}
