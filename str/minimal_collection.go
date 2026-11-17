package str

// MinimalCollection represents JSON minimal movie or episode collection: Trakt ID -> collected_at, or collected_at|source with available_on
type MinimalCollection map[string]*MinimalCollected

func (m MinimalCollection) String() string {
	return Stringify(m)
}
