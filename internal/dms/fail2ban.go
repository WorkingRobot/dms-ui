package dms

import (
	"context"
	"strconv"
	"strings"
)

// Jail is one fail2ban jail parsed from `setup fail2ban status`.
type Jail struct {
	Name            string   `json:"name"`
	CurrentlyFailed int      `json:"currentlyFailed"`
	TotalFailed     int      `json:"totalFailed"`
	CurrentlyBanned int      `json:"currentlyBanned"`
	TotalBanned     int      `json:"totalBanned"`
	BannedIPs       []string `json:"bannedIPs"`
}

// Fail2banStatus parses the concatenated per-jail status blocks.
func (c *Client) Fail2banStatus(ctx context.Context) ([]Jail, error) {
	out, err := c.setup(ctx, "fail2ban", "status")
	if err != nil {
		return nil, err
	}
	return parseFail2ban(out), nil
}

func parseFail2ban(out string) []Jail {
	var jails []Jail
	var cur *Jail
	for _, raw := range strings.Split(out, "\n") {
		ln := strings.TrimSpace(raw)
		if strings.HasPrefix(ln, "Status for the jail:") {
			if cur != nil {
				jails = append(jails, *cur)
			}
			cur = &Jail{Name: strings.TrimSpace(strings.TrimPrefix(ln, "Status for the jail:"))}
			continue
		}
		if cur == nil {
			continue
		}
		val := func() string {
			if i := strings.LastIndex(ln, ":"); i >= 0 {
				return strings.TrimSpace(ln[i+1:])
			}
			return ""
		}
		switch {
		case strings.Contains(ln, "Currently failed:"):
			cur.CurrentlyFailed, _ = strconv.Atoi(val())
		case strings.Contains(ln, "Total failed:"):
			cur.TotalFailed, _ = strconv.Atoi(val())
		case strings.Contains(ln, "Currently banned:"):
			cur.CurrentlyBanned, _ = strconv.Atoi(val())
		case strings.Contains(ln, "Total banned:"):
			cur.TotalBanned, _ = strconv.Atoi(val())
		case strings.Contains(ln, "Banned IP list:"):
			cur.BannedIPs = strings.Fields(val())
		}
	}
	if cur != nil {
		jails = append(jails, *cur)
	}
	return jails
}

// Fail2banBan bans an IP across the custom jail.
func (c *Client) Fail2banBan(ctx context.Context, ip string) (string, error) {
	return c.setup(ctx, "fail2ban", "ban", ip)
}

// Fail2banUnban removes an IP from all jails.
func (c *Client) Fail2banUnban(ctx context.Context, ip string) (string, error) {
	return c.setup(ctx, "fail2ban", "unban", ip)
}

// Fail2banLog returns the raw fail2ban log (most recent lines first capped).
func (c *Client) Fail2banLog(ctx context.Context) (string, error) {
	return c.setup(ctx, "fail2ban", "log")
}
