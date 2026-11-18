package str

// YounifyRefreshService represents JSON body of a streaming service refresh request that names the service
type YounifyRefreshService struct {
	ServiceID *string `json:"service_id,omitempty"`
	SkipSync  *bool   `json:"skip_sync,omitempty"`
}

func (y YounifyRefreshService) String() string {
	return Stringify(y)
}
