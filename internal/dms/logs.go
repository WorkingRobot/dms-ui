package dms

import (
	"context"
	"strconv"
	"strings"
)

// MailLog returns the last n lines of the mail log. DMS logs to
// /var/log/mail/mail.log (and historically /var/log/mail.log); we try both.
func (c *Client) MailLog(ctx context.Context, n int) (string, error) {
	if n <= 0 || n > 5000 {
		n = 500
	}
	r, err := c.Exec(ctx, "sh", "-c",
		"tail -n "+strconv.Itoa(n)+" /var/log/mail/mail.log 2>/dev/null || tail -n "+strconv.Itoa(n)+" /var/log/mail.log 2>/dev/null")
	if err != nil {
		return "", err
	}
	return r.Stdout, nil
}

// QueueEntry is a parsed line group from `postqueue -p`.
type QueueEntry struct {
	ID        string `json:"id"`
	Size      string `json:"size"`
	Arrival   string `json:"arrival"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Reason    string `json:"reason"`
}

// MailQueue parses `postqueue -p`. Returns an empty slice when the queue is
// empty ("Mail queue is empty").
func (c *Client) MailQueue(ctx context.Context) ([]QueueEntry, error) {
	out, err := c.run(ctx, "postqueue", "-p")
	if err != nil {
		return nil, err
	}
	if strings.Contains(out, "Mail queue is empty") {
		return []QueueEntry{}, nil
	}
	var entries []QueueEntry
	var cur *QueueEntry
	for _, raw := range strings.Split(out, "\n") {
		ln := strings.TrimRight(raw, " \t")
		if ln == "" {
			if cur != nil {
				entries = append(entries, *cur)
				cur = nil
			}
			continue
		}
		if strings.HasPrefix(ln, "-Queue ID-") || strings.HasPrefix(ln, "--") {
			continue
		}
		f := strings.Fields(ln)
		if cur == nil && len(f) >= 7 {
			cur = &QueueEntry{
				ID:      strings.TrimRight(f[0], "*!"),
				Size:    f[1],
				Arrival: strings.Join(f[2:6], " "),
				Sender:  f[6],
			}
		} else if cur != nil {
			t := strings.TrimSpace(ln)
			if strings.HasPrefix(t, "(") {
				cur.Reason = strings.Trim(t, "()")
			} else {
				cur.Recipient = t
			}
		}
	}
	if cur != nil {
		entries = append(entries, *cur)
	}
	return entries, nil
}
