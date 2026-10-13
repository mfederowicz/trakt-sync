package str

// SmartListDraft represents JSON ready-made smart list setup for a personal list that stopped updating
type SmartListDraft struct {
	List        *SmartListDraftList      `json:"list,omitempty"`
	OriginURL   *string                  `json:"origin_url,omitempty"`
	Source      *string                  `json:"source,omitempty"`
	MediaType   *string                  `json:"media_type,omitempty"`
	Filters     *SmartListFilters        `json:"filters,omitempty"`
	Unsupported []string                 `json:"unsupported,omitempty"`
	Confidence  *float64                 `json:"confidence,omitempty"`
	Converted   *SmartListDraftConverted `json:"converted,omitempty"`
}

func (s SmartListDraft) String() string {
	return Stringify(s)
}
