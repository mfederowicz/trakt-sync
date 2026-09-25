// Package uri used for url operations
package uri

// SocialActivityOptions query options for users/{id}/{type}/activities
type SocialActivityOptions struct {
	Countries string `url:"countries,omitempty"`
	Extended  string `url:"extended,omitempty"`
	Genres    string `url:"genres,omitempty"`
	Limit     int    `url:"limit,omitempty"`
	Page      int    `url:"page,omitempty"`
	Runtimes  string `url:"runtimes,omitempty"`
	Years     string `url:"years,omitempty"`
}
