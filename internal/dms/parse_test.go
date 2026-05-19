package dms

import (
	"os"
	"path/filepath"
	"testing"
)

func fixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return string(b)
}

func TestParseAccounts(t *testing.T) {
	accts := parseAccounts(fixture(t, "email_list.txt"))
	if len(accts) != 12 {
		t.Fatalf("want 12 accounts, got %d", len(accts))
	}
	a := accts[0]
	if a.Email != "admin@example.com" {
		t.Errorf("email = %q", a.Email)
	}
	if a.Used != "474M" || a.Quota != "~" || a.UsedPct != 0 {
		t.Errorf("usage parse: used=%q quota=%q pct=%d", a.Used, a.Quota, a.UsedPct)
	}
	if len(a.Aliases) != 3 || a.Aliases[0] != "sales@example.com" {
		t.Errorf("aliases = %v", a.Aliases)
	}
	if accts[8].Email != "admin@example.org" || accts[8].Used != "112K" {
		t.Errorf("account[8] = %+v", accts[8])
	}
}

func TestParseAliases(t *testing.T) {
	al := parseAliases(fixture(t, "alias_list.txt"))
	if len(al) != 3 {
		t.Fatalf("want 3 aliases, got %d (%v)", len(al), al)
	}
	if al[0].Source != "sales@example.com" || al[0].Target != "admin@example.com" {
		t.Errorf("alias[0] = %+v", al[0])
	}
}

func TestParseFail2ban(t *testing.T) {
	jails := parseFail2ban(fixture(t, "fail2ban_status.txt"))
	if len(jails) != 3 {
		t.Fatalf("want 3 jails, got %d", len(jails))
	}
	var postfix *Jail
	for i := range jails {
		if jails[i].Name == "postfix" {
			postfix = &jails[i]
		}
	}
	if postfix == nil {
		t.Fatal("postfix jail missing")
	}
	if postfix.CurrentlyBanned != 27 || postfix.TotalBanned != 59 || postfix.CurrentlyFailed != 59 {
		t.Errorf("postfix counts = %+v", postfix)
	}
	if len(postfix.BannedIPs) != 27 {
		t.Errorf("banned ips = %d", len(postfix.BannedIPs))
	}
}

func TestParseServices(t *testing.T) {
	svcs := parseServices(fixture(t, "supervisorctl_status.txt"))
	if len(svcs) != 21 {
		t.Fatalf("want 21 services, got %d", len(svcs))
	}
	found := map[string]string{}
	for _, s := range svcs {
		found[s.Name] = s.State
	}
	if found["postfix"] != "RUNNING" || found["dovecot"] != "RUNNING" {
		t.Errorf("postfix/dovecot not running: %v", found)
	}
	if found["clamav"] != "STOPPED" {
		t.Errorf("clamav state = %q", found["clamav"])
	}
}

func TestParseRspamd(t *testing.T) {
	s := parseRspamd(fixture(t, "rspamc_stat.txt"))
	if s.Scanned != 6760 || s.Spam != 102 || s.Ham != 6658 || s.Learned != 10652 {
		t.Errorf("rspamd totals = %+v", s)
	}
	if s.Actions["reject"] != 41 || s.Actions["no action"] != 6547 {
		t.Errorf("rspamd actions = %v", s.Actions)
	}
	if s.AvgScanMS != 6.164 {
		t.Errorf("avg scan = %v", s.AvgScanMS)
	}
}

func TestParseRestrictions(t *testing.T) {
	recv := parseRestrictions(fixture(t, "restrict_receive.txt"))
	if len(recv) != 9 {
		t.Fatalf("want 9 receive restrictions, got %d (%v)", len(recv), recv)
	}
	if recv[0] != "store@example.com" {
		t.Errorf("recv[0] = %q", recv[0])
	}
	// `restrict list send` here is just an INFO log line, no entries.
	send := parseRestrictions(fixture(t, "restrict_send.txt"))
	if len(send) != 0 {
		t.Errorf("want 0 send restrictions, got %v", send)
	}
}

func TestParseDKIMFilename(t *testing.T) {
	k := parseDKIMFilename("rsa-2048-mail-example.com.public.dns.txt")
	if k.KeyType != "rsa" || k.Selector != "mail" || k.Domain != "example.com" {
		t.Errorf("rsa parse = %+v", k)
	}
	k = parseDKIMFilename("ed25519-mail-example.org.public.txt")
	if k.KeyType != "ed25519" || k.Selector != "mail" || k.Domain != "example.org" {
		t.Errorf("ed25519 parse = %+v", k)
	}
}
