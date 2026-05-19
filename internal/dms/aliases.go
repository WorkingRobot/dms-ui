package dms

import (
	"context"
	"strings"
)

// Alias maps a virtual address to a recipient.
type Alias struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// ListAliases parses `setup alias list`:
//
//   - sales@example.com admin@example.com
func (c *Client) ListAliases(ctx context.Context) ([]Alias, error) {
	out, err := c.setup(ctx, "alias", "list")
	if err != nil {
		return nil, err
	}
	return parseAliases(out), nil
}

func parseAliases(out string) []Alias {
	var aliases []Alias
	for _, ln := range strings.Split(out, "\n") {
		ln = strings.TrimSpace(ln)
		ln = strings.TrimSpace(strings.TrimPrefix(ln, "*"))
		if ln == "" {
			continue
		}
		f := strings.Fields(ln)
		if len(f) >= 2 {
			aliases = append(aliases, Alias{Source: f[0], Target: f[1]})
		}
	}
	return aliases
}

// AddAlias adds a source->target virtual alias.
func (c *Client) AddAlias(ctx context.Context, source, target string) (string, error) {
	return c.setup(ctx, "alias", "add", source, target)
}

// DeleteAlias removes a single source->target pair.
func (c *Client) DeleteAlias(ctx context.Context, source, target string) (string, error) {
	return c.setup(ctx, "alias", "del", source, target)
}

// ListMasters parses `setup dovecot-master list` (one username per line).
func (c *Client) ListMasters(ctx context.Context) ([]string, error) {
	out, err := c.setup(ctx, "dovecot-master", "list")
	if err != nil {
		return nil, err
	}
	var users []string
	for _, ln := range strings.Split(out, "\n") {
		ln = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(ln), "*"))
		if ln != "" {
			users = append(users, strings.Fields(ln)[0])
		}
	}
	return users, nil
}

// AddMaster creates a Dovecot master (SASL) user.
func (c *Client) AddMaster(ctx context.Context, user, password string) (string, error) {
	return c.setup(ctx, "dovecot-master", "add", user, password)
}

// UpdateMaster changes a master user's password.
func (c *Client) UpdateMaster(ctx context.Context, user, password string) (string, error) {
	return c.setup(ctx, "dovecot-master", "update", user, password)
}

// DeleteMaster removes a master user (non-interactive via -y).
func (c *Client) DeleteMaster(ctx context.Context, user string) (string, error) {
	return c.setup(ctx, "dovecot-master", "del", "-y", user)
}
