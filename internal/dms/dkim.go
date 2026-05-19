package dms

import (
	"context"
	"path"
	"sort"
	"strings"
)

const dkimDir = "/tmp/docker-mailserver/rspamd/dkim"

// DKIMKey describes one generated DKIM key (rspamd layout:
// "<keytype>-<keysize?>-<selector>-<domain>.public.dns.txt").
type DKIMKey struct {
	Domain    string `json:"domain"`
	Selector  string `json:"selector"`
	KeyType   string `json:"keyType"`
	RecordKey string `json:"recordKey"` // "<selector>._domainkey.<domain>"
	TXTValue  string `json:"txtValue"`  // cleaned DNS value
	File      string `json:"file"`
}

// ListDKIM enumerates the rspamd DKIM public DNS files and parses each into a
// ready-to-publish TXT record.
func (c *Client) ListDKIM(ctx context.Context) ([]DKIMKey, error) {
	ls, err := c.Exec(ctx, "sh", "-c", "ls -1 "+dkimDir+"/*.public.dns.txt 2>/dev/null || true")
	if err != nil {
		return nil, err
	}
	var keys []DKIMKey
	for _, f := range strings.Fields(ls.Stdout) {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		k := parseDKIMFilename(path.Base(f))
		k.File = f
		val, err := c.run(ctx, "cat", f)
		if err == nil {
			k.TXTValue = strings.TrimSpace(strings.ReplaceAll(val, "\n", ""))
		}
		if k.Selector != "" && k.Domain != "" {
			k.RecordKey = k.Selector + "._domainkey." + k.Domain
		}
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Domain < keys[j].Domain })
	return keys, nil
}

// parseDKIMFilename: "rsa-2048-mail-example.com.public.dns.txt" or
// "ed25519-mail-example.org.public.txt". The domain may contain dots, so we
// split from the left: type [size] selector, remainder is the domain.
func parseDKIMFilename(name string) DKIMKey {
	base := name
	for _, suf := range []string{".public.dns.txt", ".public.txt", ".private.txt"} {
		base = strings.TrimSuffix(base, suf)
	}
	parts := strings.SplitN(base, "-", 2)
	k := DKIMKey{}
	if len(parts) < 2 {
		return k
	}
	k.KeyType = parts[0]
	rest := parts[1]
	if k.KeyType == "rsa" {
		// rsa-<size>-<selector>-<domain>
		p := strings.SplitN(rest, "-", 3)
		if len(p) == 3 {
			k.Selector = p[1]
			k.Domain = p[2]
		}
	} else {
		// ed25519-<selector>-<domain>
		p := strings.SplitN(rest, "-", 2)
		if len(p) == 2 {
			k.Selector = p[0]
			k.Domain = p[1]
		}
	}
	return k
}

// DKIMOptions configures key generation.
type DKIMOptions struct {
	Domain   string `json:"domain"`
	Selector string `json:"selector"` // default "mail"
	KeyType  string `json:"keyType"`  // "rsa" or "ed25519"
	KeySize  string `json:"keySize"`  // RSA only: 1024/2048/4096
	Force    bool   `json:"force"`    // overwrite/rotate existing
}

// GenerateDKIM runs `setup config dkim ...` (rspamd path) and returns the
// command output (which includes the DNS record to publish).
func (c *Client) GenerateDKIM(ctx context.Context, o DKIMOptions) (string, error) {
	args := []string{"config", "dkim"}
	if o.KeyType != "" {
		args = append(args, "keytype", o.KeyType)
	}
	if o.KeyType != "ed25519" && o.KeySize != "" {
		args = append(args, "keysize", o.KeySize)
	}
	if o.Selector != "" {
		args = append(args, "selector", o.Selector)
	}
	if o.Domain != "" {
		args = append(args, "domain", o.Domain)
	}
	if o.Force {
		args = append(args, "--force")
	}
	return c.setup(ctx, args...)
}
