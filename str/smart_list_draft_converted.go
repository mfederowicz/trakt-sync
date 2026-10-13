package str

// SmartListDraftConverted represents JSON smart list a personal list was turned into
type SmartListDraftConverted struct {
	SmartListID *int       `json:"smart_list_id,omitempty"`
	ConvertedAt *Timestamp `json:"converted_at,omitempty"`
}

func (s SmartListDraftConverted) String() string {
	return Stringify(s)
}
