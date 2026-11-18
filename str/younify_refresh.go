package str

// YounifyRefresh represents JSON body of a streaming service refresh request
type YounifyRefresh struct {
	SkipSync *bool `json:"skip_sync,omitempty"`
}

func (y YounifyRefresh) String() string {
	return Stringify(y)
}
