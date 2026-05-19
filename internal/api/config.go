package api

import "os"

// Config is resolved entirely from environment variables.
type Config struct {
	ListenAddr    string // LISTEN_ADDR (default :8080)
	MailContainer string // MAIL_CONTAINER (default mailserver)
	EmailHeader   string // PROXY_EMAIL_HEADER (default X-Forwarded-Email)
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// LoadConfig reads configuration from the environment. dms-ui has no built-in
// auth and is intended to run behind a trusted authenticating proxy
// (e.g. Traefik + Authentik); the operator is responsible for that exposure.
func LoadConfig() (*Config, error) {
	return &Config{
		ListenAddr:    env("LISTEN_ADDR", ":8080"),
		MailContainer: env("MAIL_CONTAINER", "mailserver"),
		EmailHeader:   env("PROXY_EMAIL_HEADER", "X-Forwarded-Email"),
	}, nil
}
