// Package str used for structs
package str

// Metadata represents JSON media metadata object
type Metadata struct {
	MediaType     *string `json:"media_type,omitempty"`
	Resolution    *string `json:"resolution,omitempty"`
	Hdr           *string `json:"hdr,omitempty"`
	Audio         *string `json:"audio,omitempty"`
	AudioChannels *string `json:"audio_channels,omitempty"`
	ThreeD        *bool   `json:"3d,omitempty"`
}

func (m Metadata) String() string {
	return Stringify(m)
}
