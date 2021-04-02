package umbrella

import "net/http"

const (
	chromeUserAgent  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.XX (KHTML, like Gecko) Chrome/84.0.XXXX.XX Safari/537.XX"
	firefoxUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Gecko/XXXXXXXX Firefox/76.XX"
)

// HoneypotUserAgents ...
var HoneypotUserAgents = []string{
	"0xSCANNER",
	"20010801",
	"AhrefsBot",
	"Alprazolam",
	"BLEXBot",
	"BOT for JCE",
	"Baiduspider",
	"Gecko/20100115",
	"Gemini",
	"Hakai",
	"Hello",
	"Indy Library",
	"Indy-Library",
	"JDatabaseDriverMysqli",
	"Jorgee",
	"LMAO",
	"MJ12bot",
	"Mappy",
	"Morfeus",
	"NYU",
	"Nessus",
	"Nikto",
	"OpenVAS",
	"Ronin",
	"SemrushBot",
	"Shinka",
	"WPScan",
	"ZmEu",
	"aiohttp",
	"masscan",
	"muhstik",
	"sqlmap",
	"sysscan",
	"union select",
	"yandex",
	"zgrab",
}

// AllowUserAgent is middleware that allows a request only if any of
// the specified strings is included in the User-Agent header.
// Returns 403 Forbidden status if the request has a user-agent that is not allowed.
func AllowUserAgent(userAgents ...string) func(http.Handler) http.Handler {
	return AllowHTTPHeader(http.StatusForbidden, "User-Agent", userAgents...)
}

// DisallowUserAgent is middleware that disallow a request only if any of
// the specified strings is included in the User-Agent header.
// Returns 403 Forbidden status if the request has a user-agent that is not allowed.
func DisallowUserAgent(userAgents ...string) func(http.Handler) http.Handler {
	return DisallowHTTPHeader(http.StatusForbidden, "User-Agent", userAgents...)
}
