package str

import "net/http"

type Response struct {
	*http.Response
	// Explicitly specify the Rate type so Rate's String() receiver doesn't
	// propagate to Response.
	Rate Rate
}
