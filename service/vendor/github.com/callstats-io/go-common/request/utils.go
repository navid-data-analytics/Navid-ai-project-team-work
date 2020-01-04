package request

import (
	"net/http"
	"strings"
)

// IPFromAddr returns the IP part of an inet addr
func IPFromAddr(addr string) string {
	// check for IPv6
	if strings.Count(addr, ":") > 1 {
		ip := addr
		// assume also has suffix if it has prefix for IPv6
		if strings.HasPrefix(addr, "[") {
			ip = ip[1:strings.Index(ip, "]")]
		}
		// otherwise assume only IP present and check for zone
		zoneIdx := strings.LastIndex(ip, "%")
		if zoneIdx != -1 {
			ip = ip[:zoneIdx]
		}
		return ip
	}

	// IPv4
	val := strings.LastIndex(addr, ":")
	if val == -1 {
		return addr // assume valid IP without port
	}
	ip := addr[:val]
	return ip
}

// IPFromRequest returns the IP address from a HTTP request.
// It first looks at x-csio-client-ip header and uses it if present. Otherwise it returns the remote IP address.
func IPFromRequest(r *http.Request) string {
	headerIP := r.Header.Get(CsioClientIPHeader)
	if headerIP != "" {
		return IPFromAddr(headerIP)
	}
	return IPFromAddr(r.RemoteAddr)
}
