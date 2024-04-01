package structs

func AllFactionRelations() []FactionRelation {
	all := []FactionRelation{}
	for _, r := range FactionRelation_value {
		all = append(all, FactionRelation(r))
	}
	return all
}
