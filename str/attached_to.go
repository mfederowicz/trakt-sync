package str

// AttachedTo represents JSON notes object
type AttachedTo struct {
	ID          *int       `json:"id,omitempty"`
	Type        *string    `json:"type,omitempty"`
	WatchedAt   *Timestamp `json:"watched_at,omitempty"`
	CollectedAt *Timestamp `json:"collected_at,omitempty"`
	RatedAt     *Timestamp `json:"rated_at,omitempty"`
	Rating      *int       `json:"rating,omitempty"`
}

func (a AttachedTo) String() string {
	return Stringify(a)
}
