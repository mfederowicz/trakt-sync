/*
Package uri holds the query options of the Trakt API and encodes them into request URLs.

ListOptions covers the common page, limit and extended parameters; endpoint specific filters have their
own option structs. Fields use url:"name,omitempty" tags, and AddQuery adds the set fields to a path:

	path, err := uri.AddQuery("movies/trending", &uri.ListOptions{Page: 1, Limit: 10})
	// movies/trending?limit=10&page=1

Nil options add no query.
*/
package uri
