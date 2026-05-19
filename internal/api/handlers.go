package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/WorkingRobot/dms-ui/internal/dms"
)

// okOut wraps a (commandOutput, error) DMS call into the standard handler
// response, surfacing the `setup` stdout so the UI can show it in a toast.
func okOut(out string, err error) (any, error) {
	if err != nil {
		return nil, err
	}
	return map[string]string{"status": "ok", "output": strings.TrimSpace(out)}, nil
}

// param returns a URL-decoded path parameter (emails arrive percent-encoded).
func param(r *http.Request, key string) string {
	v := chi.URLParam(r, key)
	if dec, err := url.PathUnescape(v); err == nil {
		return dec
	}
	return v
}

// ---- Dashboard ----

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	out := map[string]any{
		"container": s.dms.ContainerName(),
		"version":   s.dms.Version(ctx),
	}
	if svcs, err := s.dms.Services(ctx); err == nil {
		out["services"] = svcs
	}
	if st, err := s.dms.RspamdStats(ctx); err == nil {
		out["rspamd"] = st
	}
	if accts, err := s.dms.ListAccounts(ctx); err == nil {
		out["accountCount"] = len(accts)
	}
	if al, err := s.dms.ListAliases(ctx); err == nil {
		out["aliasCount"] = len(al)
	}
	if q, err := s.dms.MailQueue(ctx); err == nil {
		out["queueLength"] = len(q)
	}
	writeJSON(w, http.StatusOK, out)
}

// ---- Accounts ----

func (s *Server) listAccounts(r *http.Request) (any, error) {
	return s.dms.ListAccounts(r.Context())
}

func (s *Server) accountInfo(r *http.Request) (any, error) {
	return s.dms.AccountInfo(r.Context(), param(r, "email"))
}

type accountReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) addAccount(r *http.Request) (any, error) {
	var b accountReq
	if err := decode(r, &b); err != nil {
		return nil, err
	}
	if b.Email == "" || b.Password == "" {
		return nil, badRequest("email and password are required")
	}
	return okOut(s.dms.AddAccount(r.Context(), b.Email, b.Password))
}

func (s *Server) updateAccount(r *http.Request) (any, error) {
	var b struct {
		Password string `json:"password"`
	}
	if err := decode(r, &b); err != nil {
		return nil, err
	}
	if b.Password == "" {
		return nil, badRequest("password is required")
	}
	return okOut(s.dms.UpdateAccount(r.Context(), param(r, "email"), b.Password))
}

func (s *Server) deleteAccount(r *http.Request) (any, error) {
	purge := r.URL.Query().Get("purge") == "true"
	return okOut(s.dms.DeleteAccount(r.Context(), param(r, "email"), purge))
}

func (s *Server) setQuota(r *http.Request) (any, error) {
	var b struct {
		Quota string `json:"quota"`
	}
	if err := decode(r, &b); err != nil {
		return nil, err
	}
	if b.Quota == "" {
		return nil, badRequest("quota is required (e.g. 5G, 512M, or 0)")
	}
	return okOut(s.dms.SetQuota(r.Context(), param(r, "email"), b.Quota))
}

func (s *Server) delQuota(r *http.Request) (any, error) {
	return okOut(s.dms.DeleteQuota(r.Context(), param(r, "email")))
}

// ---- Restrictions ----

func (s *Server) listRestrictions(r *http.Request) (any, error) {
	return s.dms.ListRestrictions(r.Context(), param(r, "dir"))
}

type restrictReq struct {
	Email string `json:"email"`
}

func (s *Server) addRestriction(r *http.Request) (any, error) {
	var b restrictReq
	if err := decode(r, &b); err != nil {
		return nil, err
	}
	return okOut(s.dms.SetRestriction(r.Context(), "add", param(r, "dir"), b.Email))
}

func (s *Server) delRestriction(r *http.Request) (any, error) {
	email := r.URL.Query().Get("email")
	if email == "" {
		return nil, badRequest("email query param is required")
	}
	return okOut(s.dms.SetRestriction(r.Context(), "del", param(r, "dir"), email))
}

// ---- Aliases ----

func (s *Server) listAliases(r *http.Request) (any, error) {
	return s.dms.ListAliases(r.Context())
}

func (s *Server) addAlias(r *http.Request) (any, error) {
	var b dms.Alias
	if err := decode(r, &b); err != nil {
		return nil, err
	}
	if b.Source == "" || b.Target == "" {
		return nil, badRequest("source and target are required")
	}
	return okOut(s.dms.AddAlias(r.Context(), b.Source, b.Target))
}

func (s *Server) delAlias(r *http.Request) (any, error) {
	src := r.URL.Query().Get("source")
	tgt := r.URL.Query().Get("target")
	if src == "" || tgt == "" {
		return nil, badRequest("source and target query params are required")
	}
	return okOut(s.dms.DeleteAlias(r.Context(), src, tgt))
}

// ---- Dovecot masters ----

func (s *Server) listMasters(r *http.Request) (any, error) {
	return s.dms.ListMasters(r.Context())
}

type masterReq struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func (s *Server) addMaster(r *http.Request) (any, error) {
	var b masterReq
	if err := decode(r, &b); err != nil {
		return nil, err
	}
	if b.User == "" || b.Password == "" {
		return nil, badRequest("user and password are required")
	}
	return okOut(s.dms.AddMaster(r.Context(), b.User, b.Password))
}

func (s *Server) updateMaster(r *http.Request) (any, error) {
	var b struct {
		Password string `json:"password"`
	}
	if err := decode(r, &b); err != nil {
		return nil, err
	}
	return okOut(s.dms.UpdateMaster(r.Context(), param(r, "user"), b.Password))
}

func (s *Server) delMaster(r *http.Request) (any, error) {
	return okOut(s.dms.DeleteMaster(r.Context(), param(r, "user")))
}

// ---- DKIM ----

func (s *Server) listDKIM(r *http.Request) (any, error) {
	return s.dms.ListDKIM(r.Context())
}

func (s *Server) genDKIM(r *http.Request) (any, error) {
	var o dms.DKIMOptions
	if err := decode(r, &o); err != nil {
		return nil, err
	}
	if o.Domain == "" {
		return nil, badRequest("domain is required")
	}
	out, err := s.dms.GenerateDKIM(r.Context(), o)
	if err != nil {
		return nil, err
	}
	return map[string]string{"output": out}, nil
}

// ---- Relay ----

func (s *Server) relayInfo(r *http.Request) (any, error) {
	return s.dms.RelayInfo(r.Context())
}

func (s *Server) addRelayAuth(r *http.Request) (any, error) {
	var b struct {
		Domain, User, Password string
	}
	if err := decode(r, &b); err != nil {
		return nil, err
	}
	return okOut(s.dms.AddRelayAuth(r.Context(), b.Domain, b.User, b.Password))
}

func (s *Server) addRelayDomain(r *http.Request) (any, error) {
	var b struct {
		Domain, Host, Port string
	}
	if err := decode(r, &b); err != nil {
		return nil, err
	}
	return okOut(s.dms.AddRelayDomain(r.Context(), b.Domain, b.Host, b.Port))
}

func (s *Server) excludeRelay(r *http.Request) (any, error) {
	var b struct {
		Domain string
	}
	if err := decode(r, &b); err != nil {
		return nil, err
	}
	return okOut(s.dms.ExcludeRelayDomain(r.Context(), b.Domain))
}

// ---- Fail2ban ----

func (s *Server) fail2banStatus(r *http.Request) (any, error) {
	return s.dms.Fail2banStatus(r.Context())
}

type ipReq struct {
	IP string `json:"ip"`
}

func (s *Server) fail2banBan(r *http.Request) (any, error) {
	var b ipReq
	if err := decode(r, &b); err != nil {
		return nil, err
	}
	return okOut(s.dms.Fail2banBan(r.Context(), b.IP))
}

func (s *Server) fail2banUnban(r *http.Request) (any, error) {
	var b ipReq
	if err := decode(r, &b); err != nil {
		return nil, err
	}
	return okOut(s.dms.Fail2banUnban(r.Context(), b.IP))
}

func (s *Server) fail2banLog(r *http.Request) (any, error) {
	out, err := s.dms.Fail2banLog(r.Context())
	if err != nil {
		return nil, err
	}
	return map[string]string{"log": out}, nil
}

// ---- Logs / queue ----

func (s *Server) mailLog(r *http.Request) (any, error) {
	out, err := s.dms.MailLog(r.Context(), qInt(r, "lines", 500))
	if err != nil {
		return nil, err
	}
	return map[string]string{"log": out}, nil
}

func (s *Server) mailQueue(r *http.Request) (any, error) {
	return s.dms.MailQueue(r.Context())
}

// ---- Settings (read-only server env) ----

func (s *Server) settings(r *http.Request) (any, error) {
	env, err := s.dms.Settings(r.Context())
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"container": s.dms.ContainerName(),
		"version":   s.dms.Version(r.Context()),
		"env":       env,
	}, nil
}
