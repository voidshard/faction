package config

// Faith represents a single faith used in demographics
type Faith struct {
	Distribution

	ReligionID string

	Occurs float64

	IsMonotheistic bool
}
