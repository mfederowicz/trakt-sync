// Package str used for structs
package str

// Crew represents JSON crew object
type Crew struct {
	Production       *[]Job `json:"production,omitempty"`
	Art              *[]Job `json:"art,omitempty"`
	Crew             *[]Job `json:"crew,omitempty"`
	Sound            *[]Job `json:"sound,omitempty"`
	CostumeAndMakeup *[]Job `json:"costume & make-up,omitempty"`
	Writing          *[]Job `json:"writing,omitempty"`
	Editing          *[]Job `json:"editing,omitempty"`
	VisualEffects    *[]Job `json:"visual effects,omitempty"`
	Camera           *[]Job `json:"camera,omitempty"`
	Directing        *[]Job `json:"directing,omitempty"`
	Lighting         *[]Job `json:"lighting,omitempty"`
}

func (c Crew) String() string {
	return Stringify(c)
}
