package str

// SmartListDraftList represents JSON personal list a smart list draft belongs to; a single draft sends only ids
type SmartListDraftList struct {
	IDs        *IDs     `json:"ids,omitempty"`
	Name       *string  `json:"name,omitempty"`
	ItemCount  *int     `json:"item_count,omitempty"`
	Confidence *float64 `json:"confidence,omitempty"`
}

func (s SmartListDraftList) String() string {
	return Stringify(s)
}
