package umbrella

import "net/http"

// ContentSniffing adds a header for Content-Type sniffing
// vulnerability countermeasures.
func ContentSniffing() func(http.Handler) http.Handler {
	return ResponseHeader(ContentSniffingHeaderFunc())
}

// ContentSniffingHeaderFunc returns a HeaderFunc for Content-Type
// sniffing vulnerability countermeasure.
func ContentSniffingHeaderFunc() HeaderFunc {
	return AddHeaderFunc("X-Content-Type-Options", "nosniff")
}
