# traefik-real-ip-plugin

[![Tag Badge]][Tag] [![Go Version Badge]][Go Version] [![Go Report Card Badge]][Go Report Card]

Traefik plugin to retrieve client IPs. Supports retrieving the IP from and writing the result to arbitrary headers.

If a list or a count of trusted proxies is specified, the header is parsed with the same format as the `X-Forwarded-For` header (`1.2.3.4, 1.2.3.5, ...`). Otherwise, the entire header is parsed as an IP.

The first valid IP retrieved is written to the configured destination headers.

## Traefik static configuration

```yaml
experimental:
  plugins:
    realip:
      moduleName: github.com/Desuuuu/traefik-real-ip-plugin
      version: v1.0.1
```

## Dynamic configuration

### Trusted CIDRs

```yaml
http:
  middlewares:
    realip:
      plugin:
        realip:
          retrievers:
            - header: X-Forwarded-For
              proxyCIDRs:
                - "203.0.113.195/24"
          headers:
            - X-Real-IP
```

### Trusted count

```yaml
http:
  middlewares:
    realip:
      plugin:
        realip:
          retrievers:
            - header: X-Forwarded-For
              proxyCount: 1
          headers:
            - X-Real-IP
```

[Tag]: https://github.com/Desuuuu/traefik-real-ip-plugin/tags
[Tag Badge]: https://img.shields.io/github/v/tag/Desuuuu/traefik-real-ip-plugin?sort=semver
[Go Version]: /go.mod
[Go Version Badge]: https://img.shields.io/github/go-mod/go-version/Desuuuu/traefik-real-ip-plugin
[Go Report Card]: https://goreportcard.com/report/github.com/Desuuuu/traefik-real-ip-plugin
[Go Report Card Badge]: https://goreportcard.com/badge/github.com/Desuuuu/traefik-real-ip-plugin
