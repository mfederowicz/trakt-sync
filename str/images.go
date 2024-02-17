package str

type Avatar struct {
	Full *string `json:"full,omitempty"`
}

type Images struct {
	Avatar *Avatar `json:"avatar,omitempty"`
}
