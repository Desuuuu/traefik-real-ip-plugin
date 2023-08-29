package traefik_real_ip_plugin

import (
	"net"
	"net/http"
	"strings"
)

type Retriever interface {
	Retrieve(http.Header) net.IP
}

type HeaderRetriever struct {
	Header string
}

func (r *HeaderRetriever) Retrieve(headers http.Header) net.IP {
	for _, value := range headers.Values(r.Header) {
		if value == "" {
			continue
		}

		if ip := parseIP(value); ip != nil {
			return ip
		}
	}

	return nil
}

type ProxyCountRetriever struct {
	Header string
	Count  int
}

func (r *ProxyCountRetriever) Retrieve(headers http.Header) net.IP {
	if r.Count < 1 {
		return nil
	}

	for _, value := range headers.Values(r.Header) {
		if value == "" {
			continue
		}

		list := strings.Split(value, ",")

		i := len(list) - r.Count
		if i < 0 {
			continue
		}

		if ip := parseIP(list[i]); ip != nil {
			return ip
		}
	}

	return nil
}

type ProxyCIDRRetriever struct {
	Header string
	CIDRs  []*net.IPNet
}

func (r *ProxyCIDRRetriever) Retrieve(headers http.Header) net.IP {
	for _, value := range headers.Values(r.Header) {
		if value == "" {
			continue
		}

		list := strings.Split(value, ",")

		for i := len(list) - 1; i >= 0; i-- {
			ip := parseIP(list[i])
			if ip == nil {
				break
			}

			isProxy := false

			for _, cidr := range r.CIDRs {
				if cidr.Contains(ip) {
					isProxy = true
					break
				}
			}

			if !isProxy {
				return ip
			}
		}
	}

	return nil
}

func parseIP(str string) net.IP {
	str = strings.TrimSpace(str)

	if host, port, err := net.SplitHostPort(str); err == nil {
		p, err := strconv.ParseInt(port, 10, 64)
		if err == nil && p >= 0 && p <= 65535 {
			str = host
		}
	}

	return net.ParseIP(str)
}
