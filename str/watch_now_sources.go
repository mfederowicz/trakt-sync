// Package str used for structs
package str

// WatchNowSources represents JSON watch now sources for one country
type WatchNowSources struct {
	Cable          []*WatchNowService `json:"cable,omitempty"`
	Free           []*WatchNowService `json:"free,omitempty"`
	Cinema         []*WatchNowService `json:"cinema,omitempty"`
	Subscription   []*WatchNowService `json:"subscription,omitempty"`
	Purchase       []*WatchNowService `json:"purchase,omitempty"`
	StreamingRanks *WatchNowRank      `json:"streaming_ranks,omitempty"`
}

func (w WatchNowSources) String() string {
	return Stringify(w)
}
