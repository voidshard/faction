package structs

func (g *Government) ObjectID() string {
	return g.ID
}

// NewLaws returns a new Laws struct with all maps initialised.
func NewLaws() *Laws {
	return &Laws{
		Factions:    map[string]bool{},
		Actions:     map[string]bool{},
		Commodities: map[string]bool{},
		Research:    map[string]bool{},
	}
}
