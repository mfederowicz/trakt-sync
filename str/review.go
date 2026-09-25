// Package str used for structs
package str

// Review represents JSON month / year in review object; streaming_services is sent for a month only
type Review struct {
	Stats             *ReviewStats             `json:"stats,omitempty"`
	Images            *ReviewImages            `json:"images,omitempty"`
	FirstWatched      *ReviewWatchedItem       `json:"first_watched,omitempty"`
	LastWatched       *ReviewWatchedItem       `json:"last_watched,omitempty"`
	Countries         *ReviewCountries         `json:"countries,omitempty"`
	Trends            *ReviewTrends            `json:"trends,omitempty"`
	Thanks            *ReviewThanks            `json:"thanks,omitempty"`
	StreamingServices *ReviewStreamingServices `json:"streaming_services,omitempty"`
}

func (r Review) String() string {
	return Stringify(r)
}
