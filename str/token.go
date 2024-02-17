package str

import "time"

type TokenInterface interface {
	Expired() bool
	ExpiritySeconds() int
	ExpirationPoint() time.Time
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	CreatedAt    int    `json:"created_at"`
}

func (t *Token) Expired() bool {

	return time.Now().After(t.ExpirationPoint())

}

func (t *Token) ExpiritySeconds() int {

	return int(time.Until(t.ExpirationPoint()).Seconds())
}

func (t *Token) ExpirationPoint() time.Time {

	return time.Unix(int64(t.CreatedAt), 0).Add(time.Second * time.Duration(t.ExpiresIn))

}
