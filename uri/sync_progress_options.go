// Package uri used for url operations
package uri

// SyncProgressOptions query options for sync/progress/up_next and sync/progress/watched
type SyncProgressOptions struct {
	Extended         string `url:"extended,omitempty"`
	HideCompleted    bool   `url:"hide_completed,omitempty"`
	HideNotCompleted bool   `url:"hide_not_completed,omitempty"`
	IncludeStats     bool   `url:"include_stats,omitempty"`
	LifetimeStats    bool   `url:"lifetime_stats,omitempty"`
	Limit            int    `url:"limit,omitempty"`
	OnlyRewatching   bool   `url:"only_rewatching,omitempty"`
	Page             int    `url:"page,omitempty"`
	SortBy           string `url:"sort_by,omitempty"`
	SortHow          string `url:"sort_how,omitempty"`
}
