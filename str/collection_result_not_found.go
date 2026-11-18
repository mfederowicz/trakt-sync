package str

// CollectionResultNotFound represents JSON not found object
type CollectionResultNotFound struct {
	Movies   *[]ExportlistItem `json:"movies,omitempty"`
	Shows    *[]Show           `json:"shows,omitempty"`
	Seasons  *[]Show           `json:"seasons,omitempty"`
	Episodes *[]Show           `json:"episodes,omitempty"`
}

func (c CollectionResultNotFound) String() string {
	return Stringify(c)
}
