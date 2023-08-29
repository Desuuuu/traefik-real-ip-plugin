package traefik_real_ip_plugin

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type RetrieverConfig struct {
	Header        string   `json:"header"`
	Depth         int      `json:"depth,omitempty"`
	ExcludedCIDRs []string `json:"excludedCIDRs,omitempty"`
}

type Config struct {
	Retrievers []RetrieverConfig `json:"retrievers,omitempty"`
	Headers    []string          `json:"headers,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		Retrievers: nil,
		Headers:    nil,
	}
}

type RealIP struct {
	next       http.Handler
	name       string
	retrievers []Retriever
	headers    []string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config == nil {
		return nil, errors.New("invalid config")
	}

	ri := &RealIP{
		next:    next,
		name:    name,
		headers: append([]string{}, config.Headers...),
	}

	for _, rc := range config.Retrievers {
		if rc.Header == "" {
			continue
		}

		switch {
		case rc.Depth > 0:
			ri.retrievers = append(ri.retrievers, &DepthRetriever{
				Header: rc.Header,
				Depth:  rc.Depth,
			})
		case len(rc.ExcludedCIDRs) > 0:
			cidrs := make([]*net.IPNet, 0, len(rc.ExcludedCIDRs))

			for _, c := range rc.ExcludedCIDRs {
				_, cidr, err := net.ParseCIDR(c)
				if err != nil {
					return nil, fmt.Errorf("invalid CIDR: %w", err)
				}

				cidrs = append(cidrs, cidr)
			}

			ri.retrievers = append(ri.retrievers, &ExcludedCIDRRetriever{
				Header: rc.Header,
				CIDRs:  cidrs,
			})
		default:
			ri.retrievers = append(ri.retrievers, &HeaderRetriever{
				Header: rc.Header,
			})
		}
	}

	return ri, nil
}

func (ri *RealIP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, retriever := range ri.retrievers {
		if ip := retriever.Retrieve(r.Header); ip != nil {
			ipStr := ip.String()

			for _, header := range ri.headers {
				r.Header.Set(header, ipStr)
			}
			break
		}
	}

	ri.next.ServeHTTP(w, r)
}
