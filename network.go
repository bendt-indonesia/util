package util

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

func GetIP(r *http.Request) string {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip = range splitIps {
		netIP = net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	return ""
}

func IsValidURL(URL string) bool {
	_, err := url.ParseRequestURI(URL)
	if err != nil {
		return false
	}

	return true
}

func SafeGet(param url.Values, key string) string {
	if param == nil {
		return ""
	}
	vs, ok := param[key]
	if !ok || len(vs) == 0 {
		return ""
	}
	return vs[0]
}
