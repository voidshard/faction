package structs

// Relationships applies the sort.Interface to a slice of Relationships
type Relationships []*Relationship

func (r Relationships) Len() int {
	return len(r)
}

func (r Relationships) Less(i, j int) bool {
	// low -> high
	return r[i].Trust < r[j].Trust
}

func (r Relationships) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (x *Relationship) SetMemories(v []*Memory) {
	x.Memories = v
}
