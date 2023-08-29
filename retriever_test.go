package traefik_real_ip_plugin

import (
	"net"
	"net/http"
	"testing"
)

func TestHeaderRetriever_Retrieve(t *testing.T) {
	retriever := &HeaderRetriever{
		Header: "X-Forwarded-For",
	}

	t.Run("IPv4", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 192.168.0.1 ")

		expected := net.ParseIP("192.168.0.1")

		ip := retriever.Retrieve(headers)
		if ip == nil || !ip.Equal(expected) {
			t.Errorf("Expected %v, but got %v", expected, ip)
		}
	})

	t.Run("IPv4 with port", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 192.168.0.1:9999 ")

		expected := net.ParseIP("192.168.0.1")

		ip := retriever.Retrieve(headers)
		if ip == nil || !ip.Equal(expected) {
			t.Errorf("Expected %v, but got %v", expected, ip)
		}
	})

	t.Run("IPv6", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 2001:0db8:85a3:0000:0000:8a2e:0370:7334 ")

		expected := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")

		ip := retriever.Retrieve(headers)
		if ip == nil || !ip.Equal(expected) {
			t.Errorf("Expected %v, but got %v", expected, ip)
		}
	})

	t.Run("IPv6 with port", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " [2001:0db8:85a3:0000:0000:8a2e:0370:7334]:9999 ")

		expected := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")

		ip := retriever.Retrieve(headers)
		if ip == nil || !ip.Equal(expected) {
			t.Errorf("Expected %v, but got %v", expected, ip)
		}
	})

	t.Run("invalid IP", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 2.2.2 ")

		ip := retriever.Retrieve(headers)
		if ip != nil {
			t.Errorf("Expected nil, but got %v", ip)
		}
	})

	t.Run("invalid port", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 192.168.0.1:99999 ")

		ip := retriever.Retrieve(headers)
		if ip != nil {
			t.Errorf("Expected nil, but got %v", ip)
		}
	})

	t.Run("multiple IPs", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 192.168.0.1, 192.168.1.2 ")

		ip := retriever.Retrieve(headers)
		if ip != nil {
			t.Errorf("Expected nil, but got %v", ip)
		}
	})
}

func TestDepthRetriever_Retrieve(t *testing.T) {
	retriever := &DepthRetriever{
		Header: "X-Forwarded-For",
		Depth:  2,
	}

	t.Run("IPv4", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 192.168.0.1, 192.168.1.2 ")

		expected := net.ParseIP("192.168.0.1")

		ip := retriever.Retrieve(headers)
		if ip == nil || !ip.Equal(expected) {
			t.Errorf("Expected %v, but got %v", expected, ip)
		}
	})

	t.Run("IPv4 with port", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 192.168.0.1:9999, 192.168.1.2 ")

		expected := net.ParseIP("192.168.0.1")

		ip := retriever.Retrieve(headers)
		if ip == nil || !ip.Equal(expected) {
			t.Errorf("Expected %v, but got %v", expected, ip)
		}
	})

	t.Run("IPv6", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 2001:0db8:85a3:0000:0000:8a2e:0370:7334, 2001:0db8:85a3:0000:0000:8a2e:0371:7335 ")

		expected := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")

		ip := retriever.Retrieve(headers)
		if ip == nil || !ip.Equal(expected) {
			t.Errorf("Expected %v, but got %v", expected, ip)
		}
	})

	t.Run("IPv6 with port", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " [2001:0db8:85a3:0000:0000:8a2e:0370:7334]:9999, 2001:0db8:85a3:0000:0000:8a2e:0371:7335 ")

		expected := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")

		ip := retriever.Retrieve(headers)
		if ip == nil || !ip.Equal(expected) {
			t.Errorf("Expected %v, but got %v", expected, ip)
		}
	})

	t.Run("invalid IP", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 2.2.2, 192.168.1.2 ")

		ip := retriever.Retrieve(headers)
		if ip != nil {
			t.Errorf("Expected nil, but got %v", ip)
		}
	})

	t.Run("invalid port", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 192.168.0.1:99999, 192.168.1.2 ")

		ip := retriever.Retrieve(headers)
		if ip != nil {
			t.Errorf("Expected nil, but got %v", ip)
		}
	})
}

func TestExcludedCIDRRetriever_Retrieve(t *testing.T) {
	retriever := &ExcludedCIDRRetriever{
		Header: "X-Forwarded-For",
		CIDRs: []*net.IPNet{
			mustParseCIDR("192.168.1.0/24"),
			mustParseCIDR("2001:0db8:85a3:0000:0000:8a2e:0371::/112"),
		},
	}

	t.Run("IPv4", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 192.168.0.1, 192.168.1.2 ")

		expected := net.ParseIP("192.168.0.1")

		ip := retriever.Retrieve(headers)
		if ip == nil || !ip.Equal(expected) {
			t.Errorf("Expected %v, but got %v", expected, ip)
		}
	})

	t.Run("IPv4 with port", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 192.168.0.1:9999, 192.168.1.2 ")

		expected := net.ParseIP("192.168.0.1")

		ip := retriever.Retrieve(headers)
		if ip == nil || !ip.Equal(expected) {
			t.Errorf("Expected %v, but got %v", expected, ip)
		}
	})

	t.Run("IPv6", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 2001:0db8:85a3:0000:0000:8a2e:0370:7334, 2001:0db8:85a3:0000:0000:8a2e:0371:7335 ")

		expected := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")

		ip := retriever.Retrieve(headers)
		if ip == nil || !ip.Equal(expected) {
			t.Errorf("Expected %v, but got %v", expected, ip)
		}
	})

	t.Run("IPv6 with port", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " [2001:0db8:85a3:0000:0000:8a2e:0370:7334]:9999, 2001:0db8:85a3:0000:0000:8a2e:0371:7335 ")

		expected := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")

		ip := retriever.Retrieve(headers)
		if ip == nil || !ip.Equal(expected) {
			t.Errorf("Expected %v, but got %v", expected, ip)
		}
	})

	t.Run("invalid IP", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 2.2.2, 192.168.1.2 ")

		ip := retriever.Retrieve(headers)
		if ip != nil {
			t.Errorf("Expected nil, but got %v", ip)
		}
	})

	t.Run("invalid port", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Forwarded-For", " 192.168.0.1:99999, 192.168.1.2 ")

		ip := retriever.Retrieve(headers)
		if ip != nil {
			t.Errorf("Expected nil, but got %v", ip)
		}
	})
}

func mustParseCIDR(str string) *net.IPNet {
	_, cidr, err := net.ParseCIDR(str)
	if err != nil {
		panic(err)
	}

	return cidr
}
