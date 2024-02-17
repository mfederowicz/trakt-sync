package str

type SearchListItem struct {
	Type           *string         `json:"type,omitempty"`
	Score          *float32        `json:"score,omitempty"`
	Movie          *Movie          `json:"movie,omitempty"`
	Show           *Show           `json:"show,omitempty"`
	Episode        *Episode        `json:"episode,omitempty"`
	Person         *Person         `json:"person,omitempty"`
	List           *PersonalList   `json:"list,omitempty"`
	PodcastEpisode *PodcastEpisode `json:"podcast_episode,omitempty"`
	Podcast        *Podcast        `json:"podcast,omitempty"`
}

func (i SearchListItem) String() string {
	return Stringify(i)
}
