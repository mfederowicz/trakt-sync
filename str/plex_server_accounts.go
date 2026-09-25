// Package str used for structs
package str

// PlexServerAccounts represents JSON home accounts and syncable libraries of a Plex server
type PlexServerAccounts struct {
	Accounts  []*PlexAccount `json:"accounts,omitempty"`
	Libraries []*PlexLibrary `json:"libraries,omitempty"`
}

func (p PlexServerAccounts) String() string {
	return Stringify(p)
}
