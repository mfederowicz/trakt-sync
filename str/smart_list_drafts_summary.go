package str

// SmartListDraftsSummary represents JSON summary of the personal lists that are ready to become smart lists
type SmartListDraftsSummary struct {
	ConvertibleCount *int                  `json:"convertible_count,omitempty"`
	BestConfidence   *float64              `json:"best_confidence,omitempty"`
	Lists            []*SmartListDraftList `json:"lists,omitempty"`
}

func (s SmartListDraftsSummary) String() string {
	return Stringify(s)
}
