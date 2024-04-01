package structs

func (t FactionLeadership) Rulers() int64 {
	switch t {
	case FactionLeadership_Single:
		return 1
	case FactionLeadership_Dual:
		return 2
	case FactionLeadership_Triad:
		return 3
	case FactionLeadership_Council:
		return 7
	case FactionLeadership_All:
		return 1
	}
	return 0
}
