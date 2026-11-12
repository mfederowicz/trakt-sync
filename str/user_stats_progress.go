package str

// UserStatsProgress represents JSON show progress counts in user stats
type UserStatsProgress struct {
	Started  *int `json:"started,omitempty"`
	Finished *int `json:"finished,omitempty"`
	Dropped  *int `json:"dropped,omitempty"`
}

func (u UserStatsProgress) String() string {
	return Stringify(u)
}
