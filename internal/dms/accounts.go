package dms

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Account is one mailbox as reported by `setup email list`, optionally
// enriched with doveadm quota/message detail.
type Account struct {
	Email     string   `json:"email"`
	Used      string   `json:"used"`              // human, e.g. "474M" or "0"
	Quota     string   `json:"quota"`             // human, "~" means unlimited
	UsedPct   int      `json:"usedPct"`           // 0-100
	Aliases   []string `json:"aliases,omitempty"` // virtual aliases pointing here
	UsedBytes int64    `json:"usedBytes"`         // from doveadm, -1 if unknown
	Messages  int64    `json:"messages"`          // from doveadm, -1 if unknown
}

// emailListLine: "* admin@example.com ( 474M / ~ ) [0%]"
// aliasLine:     "    [ aliases -> a@x, b@x ]"

// ListAccounts parses `setup email list`.
func (c *Client) ListAccounts(ctx context.Context) ([]Account, error) {
	out, err := c.setup(ctx, "email", "list")
	if err != nil {
		return nil, err
	}
	return parseAccounts(out), nil
}

func parseAccounts(out string) []Account {
	var accts []Account
	for _, raw := range strings.Split(out, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "* ") {
			a := parseAccountLine(strings.TrimPrefix(line, "* "))
			a.UsedBytes, a.Messages = -1, -1
			accts = append(accts, a)
			continue
		}
		if strings.HasPrefix(line, "[ aliases ->") && len(accts) > 0 {
			inner := strings.TrimSuffix(strings.TrimPrefix(line, "[ aliases ->"), "]")
			for _, al := range strings.Split(inner, ",") {
				if al = strings.TrimSpace(al); al != "" {
					accts[len(accts)-1].Aliases = append(accts[len(accts)-1].Aliases, al)
				}
			}
		}
	}
	return accts
}

func parseAccountLine(s string) Account {
	a := Account{}
	// "admin@example.com ( 474M / ~ ) [0%]"
	if i := strings.Index(s, "("); i >= 0 {
		a.Email = strings.TrimSpace(s[:i])
		rest := s[i+1:]
		if j := strings.Index(rest, ")"); j >= 0 {
			usage := rest[:j]
			if parts := strings.SplitN(usage, "/", 2); len(parts) == 2 {
				a.Used = strings.TrimSpace(parts[0])
				a.Quota = strings.TrimSpace(parts[1])
			}
			if k := strings.Index(rest, "["); k >= 0 {
				pct := strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(rest[k:]), "["), "%]")
				if v, err := strconv.Atoi(strings.TrimSuffix(pct, "%")); err == nil {
					a.UsedPct = v
				}
			}
		}
	} else {
		a.Email = strings.TrimSpace(s)
	}
	return a
}

// AccountDetail augments an Account with per-folder message stats.
type AccountDetail struct {
	Account
	Folders []FolderStat `json:"folders"`
}

// FolderStat is one mailbox folder's message count and size.
type FolderStat struct {
	Name     string `json:"name"`
	Messages int64  `json:"messages"`
	VSize    int64  `json:"vsize"`
}

// AccountInfo returns quota + per-folder detail for one mailbox via doveadm.
func (c *Client) AccountInfo(ctx context.Context, email string) (*AccountDetail, error) {
	d := &AccountDetail{}
	d.Email = email
	d.UsedBytes, d.Messages = -1, -1

	if q, err := c.run(ctx, "doveadm", "quota", "get", "-u", email); err == nil {
		for _, ln := range strings.Split(q, "\n") {
			f := strings.Fields(ln)
			// "User quota STORAGE 484701 - 0"
			if len(f) >= 4 && f[2] == "STORAGE" {
				if v, e := strconv.ParseInt(f[3], 10, 64); e == nil {
					d.UsedBytes = v * 1024 // doveadm STORAGE is KiB
				}
			}
			if len(f) >= 4 && f[2] == "MESSAGE" {
				if v, e := strconv.ParseInt(f[3], 10, 64); e == nil {
					d.Messages = v
				}
			}
		}
	}

	if m, err := c.run(ctx, "doveadm", "mailbox", "status", "-u", email, "messages vsize", "*"); err == nil {
		for _, ln := range strings.Split(m, "\n") {
			ln = strings.TrimSpace(ln)
			if ln == "" {
				continue
			}
			// "INBOX messages=10806 vsize=466335898"
			f := strings.Fields(ln)
			if len(f) < 3 {
				continue
			}
			fs := FolderStat{Name: f[0]}
			for _, kv := range f[1:] {
				p := strings.SplitN(kv, "=", 2)
				if len(p) != 2 {
					continue
				}
				n, _ := strconv.ParseInt(p[1], 10, 64)
				switch p[0] {
				case "messages":
					fs.Messages = n
				case "vsize":
					fs.VSize = n
				}
			}
			d.Folders = append(d.Folders, fs)
		}
	}
	return d, nil
}

// AddAccount creates a mailbox. Password is passed positionally so `setup`
// never blocks on its interactive prompt.
func (c *Client) AddAccount(ctx context.Context, email, password string) (string, error) {
	return c.setup(ctx, "email", "add", email, password)
}

// UpdateAccount changes a mailbox password.
func (c *Client) UpdateAccount(ctx context.Context, email, password string) (string, error) {
	return c.setup(ctx, "email", "update", email, password)
}

// DeleteAccount removes a mailbox. purge controls maildir deletion; we always
// pass an explicit -y/-n so `setup` is non-interactive.
func (c *Client) DeleteAccount(ctx context.Context, email string, purge bool) (string, error) {
	flag := "-n"
	if purge {
		flag = "-y"
	}
	return c.setup(ctx, "email", "del", flag, email)
}

// SetQuota sets a per-mailbox quota ("5G", "512M", or "0" for unlimited).
func (c *Client) SetQuota(ctx context.Context, email, quota string) (string, error) {
	return c.setup(ctx, "quota", "set", email, quota)
}

// DeleteQuota removes a per-mailbox quota override.
func (c *Client) DeleteQuota(ctx context.Context, email string) (string, error) {
	return c.setup(ctx, "quota", "del", email)
}

// Restriction is a send/receive access entry.
type Restriction struct {
	Email string `json:"email"`
}

// ListRestrictions returns addresses restricted for the given direction
// ("send" or "receive").
func (c *Client) ListRestrictions(ctx context.Context, direction string) ([]string, error) {
	if direction != "send" && direction != "receive" {
		return nil, fmt.Errorf("direction must be send or receive")
	}
	out, err := c.setup(ctx, "email", "restrict", "list", direction)
	if err != nil {
		return nil, err
	}
	return parseRestrictions(out), nil
}

// restrictAddrRe matches an "address ... ACTION" entry line. `setup email
// restrict list` emits real entries like "user@dom \t\t REJECT" but also
// timestamped INFO log lines ("2026-... INFO restrict-access: Everyone is
// allowed to send mails"); only the former start with an email address.
var restrictAddrRe = regexp.MustCompile(`^[^\s@]+@[^\s@]+$`)

func parseRestrictions(out string) []string {
	var res []string
	for _, ln := range strings.Split(out, "\n") {
		ln = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(ln), "*"))
		if ln == "" || strings.Contains(ln, "INFO") {
			continue
		}
		first := strings.Fields(ln)[0]
		if restrictAddrRe.MatchString(first) {
			res = append(res, first)
		}
	}
	return res
}

// SetRestriction adds or removes a send/receive restriction.
func (c *Client) SetRestriction(ctx context.Context, action, direction, email string) (string, error) {
	if action != "add" && action != "del" {
		return "", fmt.Errorf("action must be add or del")
	}
	return c.setup(ctx, "email", "restrict", action, direction, email)
}
