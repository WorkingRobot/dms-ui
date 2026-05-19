package dms

import "context"

// AddRelayAuth stores SASL credentials for a relay domain
// (postfix-sasl-password.cf).
func (c *Client) AddRelayAuth(ctx context.Context, domain, user, password string) (string, error) {
	return c.setup(ctx, "relay", "add-auth", domain, user, password)
}

// AddRelayDomain routes a domain's outbound mail through host[:port]
// (postfix-relaymap.cf).
func (c *Client) AddRelayDomain(ctx context.Context, domain, host, port string) (string, error) {
	args := []string{"relay", "add-domain", domain, host}
	if port != "" {
		args = append(args, port)
	}
	return c.setup(ctx, args...)
}

// ExcludeRelayDomain opts a domain out of the global RELAY_HOST.
func (c *Client) ExcludeRelayDomain(ctx context.Context, domain string) (string, error) {
	return c.setup(ctx, "relay", "exclude-domain", domain)
}

// RelayConfig is the read-only view of relay-related files. DMS has no
// `relay list`, so we read the source maps directly.
type RelayConfig struct {
	RelayMap     string `json:"relayMap"`     // postfix-relaymap.cf
	SASLPassword string `json:"saslPassword"` // postfix-sasl-password.cf (masked)
}

// RelayInfo reads the relay map and (masked) SASL password file.
func (c *Client) RelayInfo(ctx context.Context) (*RelayConfig, error) {
	rc := &RelayConfig{}
	if m, err := c.Exec(ctx, "sh", "-c", "cat /tmp/docker-mailserver/postfix-relaymap.cf 2>/dev/null"); err == nil {
		rc.RelayMap = m.Stdout
	}
	// Mask the password column (domain user:pass -> domain user:****).
	if s, err := c.Exec(ctx, "sh", "-c",
		`sed -E 's/(:)[^[:space:]]+$/\1****/' /tmp/docker-mailserver/postfix-sasl-password.cf 2>/dev/null`); err == nil {
		rc.SASLPassword = s.Stdout
	}
	return rc, nil
}
