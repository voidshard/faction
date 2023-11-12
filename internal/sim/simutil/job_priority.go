package simutil

import (
	"sort"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

type JobPriorityList struct {
	bucket   map[string][]*structs.Job // jobs by some key
	priority []*structs.Job            // job IDs in order of priority
}

func NewJobPriorityList(jobs []*structs.Job, actions map[string]*config.Action) *JobPriorityList {
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Priority > jobs[j].Priority
	})

	me := &JobPriorityList{
		bucket:   map[string][]*structs.Job{},
		priority: jobs,
	}

	for _, j := range jobs { // nb. we're iterating by priority order here
		cfg, ok := actions[j.Action]
		if !ok {
			continue
		}

		// by job target area
		byArea, ok := me.bucket[j.TargetAreaID]
		if !ok {
			byArea = []*structs.Job{}
		}
		me.bucket[j.TargetAreaID] = append(byArea, j)

		// by profession
		for prof := range cfg.ProfessionWeights {
			byProf, ok := me.bucket[prof]
			if !ok {
				byProf = []*structs.Job{}
			}
			me.bucket[prof] = append(byProf, j)
		}
	}

	return me
}

func (jp *JobPriorityList) FindJobWithin(p *structs.Person, areas map[string]*structs.Area) *structs.Job {
	return jp.find(p, areas)
}

func (jp *JobPriorityList) FindJob(p *structs.Person) *structs.Job {
	return jp.find(p, nil)
}

func (jp *JobPriorityList) find(p *structs.Person, areas map[string]*structs.Area) *structs.Job {
	// applies restriction based on faction areas of control
	// ie. passed in from a FactionAreas() call
	validArea := func(areaID string) bool {
		if areas == nil {
			return true
		}
		_, ok := areas[areaID]
		return ok
	}

	// 1. highest priority job by area, profession, faction
	// 2. highest priority job in their area by profession
	byProf, ok := jp.bucket[p.PreferredProfession]
	if ok {
		var backup *structs.Job
		for _, j := range byProf {
			if !validArea(j.TargetAreaID) {
				continue
			}
			if j.PeopleNow < j.PeopleMax && j.TargetAreaID == p.AreaID {
				if j.SourceFactionID == p.PreferredFactionID {
					// ideally, it's by our faction
					return j
				} else if backup == nil {
					backup = j
				}
			}
		}
		if backup != nil {
			return backup
		}
	}

	// 3. highest priority job in their area
	byArea := jp.bucket[p.AreaID]
	if ok {
		for _, j := range byArea {
			if !validArea(j.TargetAreaID) {
				continue
			}
			if j.PeopleNow < j.PeopleMax {
				return j
			}
		}
	}

	// 4. highest priority job the person can do
	for _, j := range jp.priority {
		if !validArea(j.TargetAreaID) {
			continue
		}
		if j.PeopleNow < j.PeopleMax {
			return j
		}
	}

	return nil
}
