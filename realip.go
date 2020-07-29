package umbrella

import (
	"net"
	"net/http"
	"strings"
)

var cidrs []*net.IPNet

func init() {
	blocks := []string{
		"127.0.0.1/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}

	cidrs = make([]*net.IPNet, len(blocks))
	for i, block := range blocks {
		_, cidrs[i], _ = net.ParseCIDR(block)
	}
}

// RealIP is middleware that overwrites RemoteAddr of http.Request
// with X-Forwarded-For or X-Real-IP header.
// Validation of the X-Forwarded-For header is done from right to left.
func RealIP() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if ip := realIP(r); ip != "" {
				r.RemoteAddr = ip
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func realIP(r *http.Request) string {
	headers := []string{
		http.CanonicalHeaderKey("X-Forwarded-For"),
		http.CanonicalHeaderKey("X-Real-IP"),
		http.CanonicalHeaderKey("X-ProxyUser-Ip"),
	}
	for _, h := range headers {
		list := strings.Split(r.Header.Get(h), ",")
		for i := len(list) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(list[i])
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
				continue
			}
			return ip
		}
	}
	return ""
}

func isPrivateSubnet(ip net.IP) bool {
	for i := range cidrs {
		if cidrs[i].Contains(ip) {
			return true
		}
	}
	return false
}
