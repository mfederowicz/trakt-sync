package str

type Job struct {
	Job          *string   `json:"job,omitempty"`
	Jobs         *[]string `json:"jobs,omitempty"`
	EpisodeCount *int      `json:"episode_count,omitempty"`
	Show         *Show     `json:"show,omitempty"`
	Movie        *Movie    `json:"movie,omitempty"`
}

func (j Job) String() string {
	return Stringify(j)
}
