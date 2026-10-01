// Package uri_test holds the documentation examples of the uri package.
package uri_test

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/uri"
)

// AddQuery adds the set option fields as query parameters; nil options add none.
func ExampleAddQuery() {
	path, err := uri.AddQuery("movies/trending", &uri.ListOptions{Page: 2, Limit: 10, Extended: "full"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(path)

	path, _ = uri.AddQuery("movies/trending", nil)
	fmt.Println(path)
	// Output:
	// movies/trending?extended=full&limit=10&page=2
	// movies/trending
}
