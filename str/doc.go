/*
Package str holds the request and response types of the Trakt API that the trakt package sends and returns.

Optional fields are pointers with omitempty. Bool, Int, Int64 and String return a pointer to a value, for
filling those fields in a literal:

	comment := &str.Comment{Comment: str.String("Great movie!"), Spoiler: str.Bool(false)}

Response wraps the *http.Response of a request (with its pagination headers) and the parsed rate limit.
*/
package str
