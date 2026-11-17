package str

// MinimalShowCollection represents JSON minimal show collection: show Trakt ID -> season -> episode -> collected_at, or collected_at|source with available_on
type MinimalShowCollection map[string]map[string]MinimalCollection

func (m MinimalShowCollection) String() string {
	return Stringify(m)
}
