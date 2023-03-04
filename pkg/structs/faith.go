package structs

// Faith represents a single faith used in demographics
type Faith struct {
	ReligionID string

	Occurs float64

	Average int

	Deviation int

	IsMonotheistic bool
}
