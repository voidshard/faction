package structs

func AllJobStates() []JobState {
	all := []JobState{}
	for _, r := range JobState_value {
		all = append(all, JobState(r))
	}
	return all
}

func (j *Job) ObjectID() string {
	return j.ID
}
