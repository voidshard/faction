package config

type Condition string

const (
	// the always true condition
	ConditionAlways Condition = "always"

	// assertions about the source faction
	ConditionSrcFactionIsGovernment    Condition = "src_faction_is_government"
	ConditionSrcFactionIsNotGovernment Condition = "src_faction_not_government"
	ConditionSrcFactionIsReligion      Condition = "src_faction_is_religion"
	ConditionSrcFactionHasReligion     Condition = "src_faction_has_religion"
	ConditionSrcFactionIsCovert        Condition = "src_faction_is_covert"
	ConditionSrcFactionIsNotCovert     Condition = "src_faction_not_covert"

	// assertions about the organisation of the source faction
	ConditionSrcFactionStructurePyramid Condition = "src_faction_structure_pyramid"
	ConditionSrcFactionStructureLoose   Condition = "src_faction_structure_loose"
	ConditionSrcFactionStructureCell    Condition = "src_faction_structure_cell"
)
