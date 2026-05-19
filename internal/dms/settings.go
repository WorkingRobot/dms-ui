package dms

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
)

// EnvVar is one effective container environment variable.
type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// envNoise are runtime/system vars that aren't DMS configuration and would
// only clutter the settings view.
var envNoise = map[string]bool{
	"PATH": true, "HOME": true, "HOSTNAME": true, "TERM": true,
	"PWD": true, "SHLVL": true, "container": true, "LANG": true,
	"LC_ALL": true, "DEBIAN_FRONTEND": true,
}

// Settings returns the mail container's effective environment (the live DMS
// configuration). It is read-only: DMS reads these at entrypoint, so changes
// require editing the compose/mail.env and restarting the container. We
// surface every configured variable (including ones explicitly set blank,
// they are still part of the config), minus obvious system noise.
func (c *Client) Settings(ctx context.Context) ([]EnvVar, error) {
	r, err := c.docker(ctx, "inspect", "-f", "{{json .Config.Env}}", c.container)
	if err != nil {
		return nil, err
	}
	var raw []string
	if e := json.Unmarshal([]byte(strings.TrimSpace(r.Stdout)), &raw); e != nil {
		return nil, e
	}
	var out []EnvVar
	for _, kv := range raw {
		i := strings.IndexByte(kv, '=')
		if i < 0 {
			continue
		}
		k, v := kv[:i], kv[i+1:]
		if envNoise[k] {
			continue
		}
		out = append(out, EnvVar{Key: k, Value: v})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out, nil
}
