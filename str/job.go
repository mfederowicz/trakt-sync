// Package str used for structs
package str

// Job represents JSON crew member positions
type Job struct {
	Job          *string   `json:"job,omitempty"`
	Jobs         *[]string `json:"jobs,omitempty"`
	Person       *Person   `json:"person,omitempty"`
	EpisodeCount *int      `json:"episode_count,omitempty"`
	Show         *Show     `json:"show,omitempty"`
	Movie        *Movie    `json:"movie,omitempty"`
}

func (j Job) String() string {
	return Stringify(j)
}
