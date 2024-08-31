// Package str used for structs
package str

import "time"

// TokenInterface methods for tokens
type TokenInterface interface {
	Expired() bool
	ExpiritySeconds() int
	ExpirationPoint() time.Time
}

// Token represents JSON token object
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	CreatedAt    int    `json:"created_at"`
}

// Expired check if token is expired
func (t *Token) Expired() bool {

	return time.Now().After(t.ExpirationPoint())

}

// ExpiritySeconds return number of seconds to token expire
func (t *Token) ExpiritySeconds() int {

	return int(time.Until(t.ExpirationPoint()).Seconds())
}

// ExpirationPoint return point in time with nanosecond precision for token
func (t *Token) ExpirationPoint() time.Time {

	return time.Unix(int64(t.CreatedAt), 0).Add(time.Second * time.Duration(t.ExpiresIn))

}
