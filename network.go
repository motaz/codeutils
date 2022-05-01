package codeutils

import (
	"net/http"
	"strings"
)

func GetRemoteIP(r *http.Request) (ip string) {

	ip = r.RemoteAddr
	if !strings.Contains(ip, "::") {
		ip = ip[:strings.Index(ip, ":")]
	}
	if strings.HasPrefix(ip, "127.") || strings.HasPrefix(ip, "0.") ||
		strings.HasPrefix(ip, "0:") || strings.HasPrefix(r.RemoteAddr, "[::1]") {
		ipt := r.Header.Get("X-Forwarded-For")
		if ipt != "" {
			ip = ipt
		}
	}

	return
}
