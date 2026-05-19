package dms

import (
	"context"
	"strconv"
	"strings"
)

// Service is one supervisord-managed process.
type Service struct {
	Name   string `json:"name"`
	State  string `json:"state"`  // RUNNING / STOPPED / ...
	Detail string `json:"detail"` // "pid 8, uptime 11 days, ..." or "Not started"
}

// Services parses `supervisorctl status`.
func (c *Client) Services(ctx context.Context) ([]Service, error) {
	// supervisorctl status exits non-zero when any service is STOPPED, which
	// is normal for DMS, so use Exec and ignore the exit code.
	r, err := c.Exec(ctx, "supervisorctl", "status")
	if err != nil {
		return nil, err
	}
	return parseServices(r.Stdout), nil
}

func parseServices(out string) []Service {
	var svcs []Service
	for _, ln := range strings.Split(out, "\n") {
		f := strings.Fields(ln)
		if len(f) < 2 {
			continue
		}
		svcs = append(svcs, Service{
			Name:   f[0],
			State:  f[1],
			Detail: strings.TrimSpace(strings.Join(f[2:], " ")),
		})
	}
	return svcs
}

// RspamdStat is a subset of `rspamc stat` worth surfacing on the dashboard.
type RspamdStat struct {
	Scanned   int64            `json:"scanned"`
	Spam      int64            `json:"spam"`
	Ham       int64            `json:"ham"`
	Learned   int64            `json:"learned"`
	Actions   map[string]int64 `json:"actions"`
	AvgScanMS float64          `json:"avgScanMs"`
}

// RspamdStats parses `rspamc stat`.
func (c *Client) RspamdStats(ctx context.Context) (*RspamdStat, error) {
	out, err := c.run(ctx, "rspamc", "stat")
	if err != nil {
		return nil, err
	}
	return parseRspamd(out), nil
}

func parseRspamd(out string) *RspamdStat {
	s := &RspamdStat{Actions: map[string]int64{}}
	for _, ln := range strings.Split(out, "\n") {
		ln = strings.TrimSpace(ln)
		switch {
		case strings.HasPrefix(ln, "Messages scanned:"):
			s.Scanned = atoiField(ln)
		case strings.HasPrefix(ln, "Messages treated as spam:"):
			s.Spam = atoiField(ln)
		case strings.HasPrefix(ln, "Messages treated as ham:"):
			s.Ham = atoiField(ln)
		case strings.HasPrefix(ln, "Messages learned:"):
			s.Learned = atoiField(ln)
		case strings.HasPrefix(ln, "Messages with action "):
			rest := strings.TrimPrefix(ln, "Messages with action ")
			if i := strings.Index(rest, ":"); i >= 0 {
				name := strings.TrimSpace(rest[:i])
				num := strings.TrimSpace(rest[i+1:])
				if j := strings.Index(num, ","); j >= 0 {
					num = num[:j]
				}
				if v, e := strconv.ParseInt(strings.TrimSpace(num), 10, 64); e == nil {
					s.Actions[name] = v
				}
			}
		case strings.HasPrefix(ln, "Average scan time:"):
			fields := strings.Fields(ln)
			if len(fields) >= 4 {
				s.AvgScanMS, _ = strconv.ParseFloat(fields[3], 64)
			}
		}
	}
	return s
}

// atoiField pulls the first integer that follows the last colon, e.g.
// "Messages scanned: 6760" -> 6760 ; "treated as spam: 102, 1.51%" -> 102.
func atoiField(ln string) int64 {
	i := strings.LastIndex(ln, ":")
	if i < 0 {
		return 0
	}
	rest := strings.TrimSpace(ln[i+1:])
	if j := strings.IndexAny(rest, ", "); j >= 0 {
		rest = rest[:j]
	}
	v, _ := strconv.ParseInt(rest, 10, 64)
	return v
}

// Version returns the docker-mailserver version (image label).
func (c *Client) Version(ctx context.Context) string {
	r, err := c.docker(ctx, "inspect", "-f",
		`{{ index .Config.Labels "org.opencontainers.image.version" }}`, c.container)
	if err != nil || r.ExitCode != 0 {
		return ""
	}
	v := strings.TrimSpace(r.Stdout)
	if v == "<no value>" {
		return ""
	}
	return v
}
