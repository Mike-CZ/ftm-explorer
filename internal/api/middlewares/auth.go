package middlewares

import (
	"ftm-explorer/internal/auth"
	"net"
	"net/http"
	"strings"
)

// AuthMiddleware defines HTTP handler middleware for logging incoming communication through provided ILogger.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(auth.SetIpAddress(r.Context(), getIP(r))))
	})
}

func getIP(req *http.Request) string {
	// Try to get the IP address from the X-Forwarded-For header
	if fwdAddress := req.Header.Get("X-Forwarded-For"); fwdAddress != "" {
		// The header can contain multiple IP addresses, so take the first one
		split := strings.Split(fwdAddress, ",")
		if len(split) > 0 {
			return split[0]
		}
	}

	// If the above header doesn't exist, then fall back to the direct connection's remote address
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return req.RemoteAddr // This might not be a perfect IP:port format, but it's something
	}
	return ip
}
