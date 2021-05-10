package umbrella

import (
	"net/http"
	"strings"
)

// Clickjacking mitigates clickjacking attacks by limiting the display
// of iframe.
func Clickjacking(opt string) func(http.Handler) http.Handler {
	return ResponseHeader(ClickjackingHeaderFunc(opt))
}

// ClickjackingHeaderFunc returns a HeaderFunc to mitigate a
// clickjacking vulnerability.
func ClickjackingHeaderFunc(opt string) HeaderFunc {
	opt = strings.ToLower(opt)
	if opt != "deny" && opt != "sameorigin" {
		opt = "deny"
	}
	return AddHeaderFunc("X-Frame-Options", opt)
}
