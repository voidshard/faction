package structs

// ResearchTopic is something a faction can research
type ResearchTopic struct {
	// Name of the topic, must be unique amongst all topics
	Name string

	// Probability that the topic will be researched
	Probability float64

	// Profession connected to researching this topic
	Profession string

	// Ethos that this topic is connected to
	Ethos Ethos
}
