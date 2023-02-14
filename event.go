package faction

type Event struct {
	FactionSource string
	FactionTarget string

	ActionType ActionType

	Message string
}
