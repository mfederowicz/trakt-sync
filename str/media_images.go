package str

// MediaImages represents JSON images object of a movie or a show (extended info "images")
type MediaImages struct {
	Fanart   *[]string `json:"fanart,omitempty"`
	Poster   *[]string `json:"poster,omitempty"`
	Logo     *[]string `json:"logo,omitempty"`
	Clearart *[]string `json:"clearart,omitempty"`
	Banner   *[]string `json:"banner,omitempty"`
	Thumb    *[]string `json:"thumb,omitempty"`
}

func (m MediaImages) String() string {
	return Stringify(m)
}
