package structs

func (p *Person) ObjectID() string {
	return p.ID
}

// SetBirthTick updates the birthtick, and moves the natural death tick and adulthood tick
// accordingly.
func (p *Person) SetBirthTick(t int64) {
	lifespan := p.NaturalDeathTick - p.BirthTick
	childhood := p.AdulthoodTick - p.BirthTick

	p.BirthTick = t
	p.NaturalDeathTick = t + lifespan
	p.AdulthoodTick = t + childhood
}
