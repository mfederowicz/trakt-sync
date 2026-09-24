// Package str used for structs
package str

// WatchNowService represents JSON watch now service object
type WatchNowService struct {
	Source      *string            `json:"source,omitempty"`
	Link        *string            `json:"link,omitempty"`
	UHD         *bool              `json:"uhd,omitempty"`
	Currency    *string            `json:"currency,omitempty"`
	Prices      *WatchNowPrices    `json:"prices,omitempty"`
	LinkTvos    *string            `json:"link_tvos,omitempty"`
	LinkDirect  *string            `json:"link_direct,omitempty"`
	LinkAndroid *string            `json:"link_android,omitempty"`
	LinkWebos   *WatchNowWebosLink `json:"link_webos,omitempty"`
}

func (w WatchNowService) String() string {
	return Stringify(w)
}
