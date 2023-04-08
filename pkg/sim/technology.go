package sim

import (
	"github.com/voidshard/faction/pkg/structs"
)

type Technology interface {
	// Topics returns all topics that can be researched
	Topics() []*structs.ResearchTopic

	// Topic returns the topic with the given name
	Topic(name string) *structs.ResearchTopic

	// Research returns the new research level for the given job.
	// This also gives the user a chance to do something
	// with the research, such as add a new technology or some
	// other hook.
	Research(jobID string, value int) int
}
