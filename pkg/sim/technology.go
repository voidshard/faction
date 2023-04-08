package sim

import (
	"github.com/voidshard/faction/pkg/structs"
)

type Technology interface {
	// Topics returns all topics that can be researched in the given areas
	Topics(areas ...string) []*structs.ResearchTopic

	// Topic returns the topic with the given name
	Topic(name string) *structs.ResearchTopic
}
