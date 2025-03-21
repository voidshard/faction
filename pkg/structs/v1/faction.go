package v1

type Faction struct {
	Meta `json:",inline" yaml:",inline"`

	// Health holds "health" information about the faction.
	Health FactionHealth `json:"Health" yaml:"Health" validate:"dive"`

	// Bonus is the bonus this faction provides on progress towards actions.
	Bonus FactionBonus `json:"Bonus" yaml:"Bonus" validate:"dive"`

	// Headquarters is the location of the factions headquarters.
	Headquarters Headquarters `json:"Headquarters" yaml:"Headquarters" validate:"dive"`

	// Actions configures the actions this faction prefers to take & modifiers.
	Actions ActionSelection `yaml:"Actions" json:"Actions" validate:"dive"`
}

// FactionBonus holds the bonus this faction provides on progress towards actions.
type FactionBonus struct {
	// YieldByProfession is the bonus this faction provides on progress towards actions.
	// This is added to yields defined on specific actions for the given profession.
	YieldByProfession map[string]float64 `json:"YieldByProfession" yaml:"YieldByProfession" validate:"min=1,max=1000000,dive,keys,alphanum,endkeys"`

	// Secrecy is the bonus this faction provides towards 'Secrecy' rolls
	// (see Progress.Secrecy)
	Secrecy Distribution `json:"Secrecy" yaml:"Secrecy" validate:"dive"`
}

// FactionHealth holds "health" information about a faction.
//
// Each of these is essentially a "resource" that can be spent, lost or gained.
type FactionHealth struct {
	// Wealth is the amount of money (silver, or some other form of liquid wealth) this faction has.
	// A faction with no wealth may be able to borrow funds (and be in the red for a while) but
	// sustained debt will lead them to disband.
	//
	// Wealth is gained through trade, taxes, theft etc and lost through spending on actions,
	// bribes, blackmail, expansion etc
	Wealth float64 `json:"Wealth" yaml:"Wealth"`

	// Corruption present in the faction. This has various effects
	// like making operations more expensive, reducing the effectiveness of actions etc.
	//
	// Corruption is gained through enemy actions and lost through internal actions to
	// root out corruption.
	Corruption float64 `json:"Corruption" yaml:"Corruption" validate:"gte=0"`

	// Cohesion is a measure of how well the faction sticks together. A low cohesion faction
	// may split apart.
	//
	// Cohesion is lost through failures, internal strife, betrayal, enemy actions etc and
	// gained through faction successes and internal actions to increase cohesion.
	Cohesion float64 `json:"Cohesion" yaml:"Cohesion" validate:"gte=0"`

	// Secrecy is a measure of how well the faction can keep secrets. A low secrecy faction
	// may have its plans discovered by others.
	//
	// Secrecy is gained via internal actions to bolster secrecy, plug leaks and actions
	// against other factions. It is lost through failures, internal strife, betrayal and
	// enemy actions to uncover secrets.
	Secrecy float64 `json:"Secrecy" yaml:"Secrecy" validate:"gte=0"`
}

// Headquarters is the location of the factions headquarters. A faction must have a headquarters.
type Headquarters struct {
	// Area is the area the faction is based in.
	Area string `json:"Area" yaml:"Area" validate:"uuid4"`

	// Building is the building the faction is based in.
	Building string `json:"Building" yaml:"Building" validate:"uuid4"`
}

func (x *Faction) New(in interface{}) (Object, error) {
	i := &Faction{}
	err := unmarshalObject(in, i)
	i.Kind = "faction"
	return i, err
}

/*
func (x *Faction) SetAllies(v []*Relationship) {
	x.Allies = v
}

func (x *Faction) SetEnemies(v []*Relationship) {
	x.Enemies = v
}
*/
