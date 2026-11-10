// Package clientip resolves the real client address of an HTTP request,
// trusting forwarding headers only from explicitly configured proxies.
package codeutils

import (
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"sync/atomic"
)

// Resolver determines client IPs using a fixed set of trusted proxy ranges.
type Resolver struct {
	trusted []netip.Prefix
}

// NewResolver builds a Resolver from CIDRs (or bare IPs) of proxies whose
// X-Forwarded-For entries may be believed.
func NewResolver(trusted []string) (*Resolver, error) {
	r := &Resolver{}
	for _, s := range trusted {
		p, err := parseTrusted(s)
		if err != nil {
			return nil, err
		}
		r.trusted = append(r.trusted, p)
	}
	return r, nil
}

// GetRemoteIP returns the client IP as a string, or "" if it can't be determined.
// All branching lives in the helpers below.
func (r *Resolver) GetRemoteIP(req *http.Request) string {
	return addrString(r.ClientIP(req))
}

// ClientIP returns the client address and true, or the zero Addr and false.
func (r *Resolver) ClientIP(req *http.Request) (netip.Addr, bool) {
	peer, ok := peerAddr(req.RemoteAddr)
	if !ok {
		return netip.Addr{}, false
	}
	return r.resolve(peer, forwardedFor(req)), true
}

// resolve applies the trust rules: an untrusted peer is the client itself;
// otherwise walk the forwarded chain right to left past trusted proxies.
func (r *Resolver) resolve(peer netip.Addr, hops []string) netip.Addr {
	if !r.isTrusted(peer) {
		return peer
	}
	current := peer
	for i := len(hops) - 1; i >= 0; i-- {
		ip, ok := parseAddr(hops[i])
		if !ok {
			break // nothing left of a malformed entry can be trusted
		}
		current = ip
		if !r.isTrusted(ip) {
			break
		}
	}
	return current
}

func (r *Resolver) isTrusted(a netip.Addr) bool {
	for _, p := range r.trusted {
		if p.Contains(a) {
			return true
		}
	}
	return false
}

// forwardedFor collects every X-Forwarded-For entry, in order, across all header lines.
func forwardedFor(req *http.Request) []string {
	joined := strings.Join(req.Header.Values("X-Forwarded-For"), ",")
	if joined == "" {
		return nil
	}
	return strings.Split(joined, ",")
}

// peerAddr extracts the TCP peer address from RemoteAddr (with or without port).
func peerAddr(remoteAddr string) (netip.Addr, bool) {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}
	return parseAddr(host)
}

// parseAddr strips any IPv6 zone and unmaps ::ffff:a.b.c.d to plain IPv4.
func parseAddr(s string) (netip.Addr, bool) {
	a, err := netip.ParseAddr(strings.TrimSpace(s))
	if err != nil {
		return netip.Addr{}, false
	}
	return a.WithZone("").Unmap(), true
}

// parseTrusted accepts a CIDR or a bare IP.
func parseTrusted(s string) (netip.Prefix, error) {
	s = strings.TrimSpace(s)
	if p, err := netip.ParsePrefix(s); err == nil {
		return p.Masked(), nil
	}
	a, ok := parseAddr(s)
	if !ok {
		return netip.Prefix{}, fmt.Errorf("clientip: invalid trusted proxy %q", s)
	}
	return netip.PrefixFrom(a, a.BitLen()), nil
}

// addrString renders an address, or "" when it is unknown.
func addrString(a netip.Addr, ok bool) string {
	if !ok {
		return ""
	}
	return a.String()
}

// ---------------------------------------------------------------------------
// Backward-compatible package-level API
// ---------------------------------------------------------------------------

// internalRanges reproduces the old isInternal() semantics: loopback, private,
// link-local unicast and unspecified addresses are treated as trusted proxies.
var internalRanges = []string{
	"127.0.0.0/8", "::1/128", // loopback
	"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "fc00::/7", // private
	"169.254.0.0/16", "fe80::/10", // link-local unicast
	"0.0.0.0/32", "::/128", // unspecified
}

// defaultResolver is initialized by a package-level variable initializer
// (no init() function) and can be swapped at runtime via SetTrustedProxies.
var defaultResolver = newDefaultResolver()

func newDefaultResolver() *atomic.Pointer[Resolver] {
	r, err := NewResolver(internalRanges)
	if err != nil {
		panic(err) // static list; can only fail on a programming error
	}
	p := new(atomic.Pointer[Resolver])
	p.Store(r)
	return p
}

// SetTrustedProxies replaces the default trusted-proxy list used by
// GetRemoteIP. Call it once at startup. Safe for concurrent use.
// Passing nil/empty means headers are never trusted.
func SetTrustedProxies(cidrs []string) error {
	r, err := NewResolver(cidrs)
	if err != nil {
		return err
	}
	defaultResolver.Store(r)
	return nil
}

// GetRemoteIP keeps the original signature and default behavior, so existing
// callers need no changes.
//
// Steps:
//  1. Identify the TCP peer from r.RemoteAddr (port stripped, normalized).
//  2. If the peer is not a trusted proxy, it IS the client: headers ignored.
//  3. Otherwise walk X-Forwarded-For right to left; the first untrusted hop is
//     the client. If every hop is trusted (e.g. a LAN client), the leftmost
//     hop is returned.
//
// Returns "" if the peer address cannot be parsed.
func GetRemoteClientIP(r *http.Request) string {
	res := defaultResolver.Load()

	// 1. TCP peer.
	peer, ok := peerAddr(r.RemoteAddr)
	if !ok {
		return ""
	}

	// 2 & 3. Trust decision and forwarded-chain walk.
	return res.resolve(peer, forwardedFor(r)).String()
}
