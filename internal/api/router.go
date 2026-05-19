package api

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/WorkingRobot/dms-ui/internal/dms"
)

// Server wires the DMS client and config into an HTTP handler.
type Server struct {
	cfg *Config
	dms *dms.Client
}

// NewServer builds the router. staticFS is the embedded SvelteKit build.
func NewServer(cfg *Config, client *dms.Client, staticFS fs.FS) http.Handler {
	s := &Server{cfg: cfg, dms: client}

	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/whoami", s.whoami)
		r.Get("/dashboard", s.dashboard)

		r.Get("/accounts", h(s.listAccounts))
		r.Post("/accounts", h(s.addAccount))
		r.Get("/accounts/{email}", h(s.accountInfo))
		r.Put("/accounts/{email}/password", h(s.updateAccount))
		r.Delete("/accounts/{email}", h(s.deleteAccount))
		r.Put("/accounts/{email}/quota", h(s.setQuota))
		r.Delete("/accounts/{email}/quota", h(s.delQuota))

		r.Get("/restrictions/{dir}", h(s.listRestrictions))
		r.Post("/restrictions/{dir}", h(s.addRestriction))
		r.Delete("/restrictions/{dir}", h(s.delRestriction))

		r.Get("/aliases", h(s.listAliases))
		r.Post("/aliases", h(s.addAlias))
		r.Delete("/aliases", h(s.delAlias))

		r.Get("/masters", h(s.listMasters))
		r.Post("/masters", h(s.addMaster))
		r.Put("/masters/{user}/password", h(s.updateMaster))
		r.Delete("/masters/{user}", h(s.delMaster))

		r.Get("/dkim", h(s.listDKIM))
		r.Post("/dkim", h(s.genDKIM))

		r.Get("/relay", h(s.relayInfo))
		r.Post("/relay/auth", h(s.addRelayAuth))
		r.Post("/relay/domain", h(s.addRelayDomain))
		r.Post("/relay/exclude", h(s.excludeRelay))

		r.Get("/fail2ban", h(s.fail2banStatus))
		r.Post("/fail2ban/ban", h(s.fail2banBan))
		r.Post("/fail2ban/unban", h(s.fail2banUnban))
		r.Get("/fail2ban/log", h(s.fail2banLog))

		r.Get("/logs", h(s.mailLog))
		r.Get("/queue", h(s.mailQueue))

		r.Get("/settings", h(s.settings))
	})

	r.Handle("/*", spaHandler(staticFS))
	return r
}

// --- helpers ---

type apiError struct {
	status int
	msg    string
}

func (e apiError) Error() string { return e.msg }

func badRequest(msg string) error { return apiError{http.StatusBadRequest, msg} }

// h adapts a handler that returns (any, error) into an http.HandlerFunc with
// uniform JSON encoding and error mapping.
func h(fn func(*http.Request) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, err := fn(r)
		if err != nil {
			status := http.StatusInternalServerError
			if ae, ok := err.(apiError); ok {
				status = ae.status
			}
			writeJSON(w, status, map[string]string{"error": err.Error()})
			return
		}
		if v == nil {
			v = map[string]string{"status": "ok"}
		}
		writeJSON(w, http.StatusOK, v)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func decode(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return badRequest("invalid JSON body: " + err.Error())
	}
	return nil
}

func qInt(r *http.Request, key string, def int) int {
	if v := r.URL.Query().Get(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func (s *Server) whoami(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get(s.cfg.EmailHeader)
	if email == "" {
		email = r.Header.Get("X-Forwarded-User")
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"email":     email,
		"container": s.dms.ContainerName(),
	})
}

// spaHandler serves the embedded static build, falling back to index.html so
// client-side routes resolve (SPA behaviour).
func spaHandler(fsys fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(fsys))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/")
		if p == "" {
			p = "index.html"
		}
		if _, err := fs.Stat(fsys, p); err != nil {
			r2 := r.Clone(r.Context())
			r2.URL.Path = "/"
			http.ServeFileFS(w, r2, fsys, "index.html")
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}
