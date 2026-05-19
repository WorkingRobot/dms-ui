// Static metadata for docker-mailserver environment variables and supervisord
// services, reconciled against the canonical upstream mailserver.env
// (https://github.com/docker-mailserver/docker-mailserver/blob/master/mailserver.env)
// and the environment reference docs. Used to give the Environment and
// Dashboard pages context instead of raw lists. Defaults shown are the
// *effective* runtime defaults DMS applies, not the (often blank) file values.

// truthy interprets a DMS on/off env value.
export const truthy = (v) => v === '1' || v === 'true' || v === 'yes';

// Ordered groups so the Environment page renders in a sensible order.
export const ENV_GROUPS = [
	'General',
	'Security & Anti-Spam',
	'Rspamd',
	'SpamAssassin',
	'Mail Fetching',
	'Reporting & Logs',
	'Relay',
	'TLS / SSL',
	'LDAP & External Auth',
	'Other'
];

// key -> { group, desc, default? }
export const ENV_META = {
	// General
	OVERRIDE_HOSTNAME: { group: 'General', desc: 'FQDN used when the container hostname cannot be set normally.' },
	LOG_LEVEL: { group: 'General', desc: 'Verbosity of DMS startup scripts and the change detector.', default: 'info' },
	SUPERVISOR_LOGLEVEL: { group: 'General', desc: 'Supervisor logging verbosity.', default: 'warn' },
	ACCOUNT_PROVISIONER: { group: 'General', desc: 'Where user accounts come from (FILE / LDAP / OIDC).', default: 'FILE' },
	PERMIT_DOCKER: { group: 'General', desc: 'Which Docker networks Postfix trusts as internal.', default: 'none' },
	TZ: { group: 'General', desc: 'Container timezone (AREA/ZONE).' },
	DMS_VMAIL_UID: { group: 'General', desc: 'UID of the vmail user owning mail storage.', default: '5000' },
	DMS_VMAIL_GID: { group: 'General', desc: 'GID of the vmail user.', default: '5000' },
	SMTP_ONLY: { group: 'General', desc: 'Run SMTP only (disables IMAP/POP3 and Dovecot).', default: '0' },
	ENABLE_IMAP: { group: 'General', desc: 'Enable IMAP access.', default: '1' },
	ENABLE_POP3: { group: 'General', desc: 'Enable POP3 retrieval.', default: '0' },
	ENABLE_MANAGESIEVE: { group: 'General', desc: 'Enable ManageSieve (port 4190) for user-managed Sieve scripts.', default: '0' },
	DOVECOT_MAILBOX_FORMAT: { group: 'General', desc: 'On-disk mailbox format (maildir / sdbox / mdbox).', default: 'maildir' },
	POSTMASTER_ADDRESS: { group: 'General', desc: 'Address that receives administrative/bounce notices.' },
	ENABLE_UPDATE_CHECK: { group: 'General', desc: 'Daily check for newer DMS releases.', default: '1' },
	UPDATE_CHECK_INTERVAL: { group: 'General', desc: 'How often to check for updates.', default: '1d' },
	POSTFIX_INET_PROTOCOLS: { group: 'General', desc: 'IP protocols Postfix listens on (all / ipv4 / ipv6).', default: 'all' },
	DOVECOT_INET_PROTOCOLS: { group: 'General', desc: 'IP protocols Dovecot listens on (all / ipv4 / ipv6).', default: 'all' },
	NETWORK_INTERFACE: { group: 'General', desc: 'Interface services bind to when not eth0.', default: 'eth0' },
	ENABLE_QUOTAS: { group: 'General', desc: 'Per-user Dovecot mailbox quotas.', default: '1' },
	POSTFIX_MAILBOX_SIZE_LIMIT: { group: 'General', desc: 'Max mailbox size in bytes (0 = unlimited).', default: '0' },
	POSTFIX_MESSAGE_SIZE_LIMIT: { group: 'General', desc: 'Max single message size in bytes.', default: '10240000' },
	POSTFIX_DAGENT: { group: 'General', desc: 'Override the LMTP/delivery agent Postfix hands mail to.' },

	// Security & Anti-Spam
	ENABLE_FAIL2BAN: { group: 'Security & Anti-Spam', desc: 'Ban hosts after repeated auth failures.', default: '0' },
	FAIL2BAN_BLOCKTYPE: { group: 'Security & Anti-Spam', desc: 'How fail2ban rejects banned IPs (drop / reject).', default: 'drop' },
	ENABLE_CLAMAV: { group: 'Security & Anti-Spam', desc: 'ClamAV virus scanning (needs ~1 GB RAM).', default: '0' },
	CLAMAV_MESSAGE_SIZE_LIMIT: { group: 'Security & Anti-Spam', desc: 'Skip ClamAV for messages larger than this.', default: '25M' },
	ENABLE_AMAVIS: { group: 'Security & Anti-Spam', desc: 'Amavis content-filter front-end (SA/ClamAV glue).', default: '1' },
	AMAVIS_LOGLEVEL: { group: 'Security & Anti-Spam', desc: 'Amavis log verbosity (-1 to 5).', default: '0' },
	ENABLE_OPENDKIM: { group: 'Security & Anti-Spam', desc: 'OpenDKIM signing (disabled automatically when Rspamd signs).', default: '1' },
	ENABLE_OPENDMARC: { group: 'Security & Anti-Spam', desc: 'OpenDMARC policy verification.', default: '1' },
	ENABLE_POLICYD_SPF: { group: 'Security & Anti-Spam', desc: 'SPF policy checks in Postfix.', default: '1' },
	ENABLE_DNSBL: { group: 'Security & Anti-Spam', desc: 'DNS block lists in Postscreen.', default: '0' },
	ENABLE_POSTGREY: { group: 'Security & Anti-Spam', desc: 'Postgrey greylisting.', default: '0' },
	POSTGREY_DELAY: { group: 'Security & Anti-Spam', desc: 'Greylist delay in seconds.', default: '300' },
	POSTGREY_MAX_AGE: { group: 'Security & Anti-Spam', desc: 'Forget greylist entries not seen for this many days.', default: '35' },
	POSTGREY_AUTO_WHITELIST_CLIENTS: { group: 'Security & Anti-Spam', desc: 'Auto-whitelist a host after this many successful deliveries (0 disables).', default: '5' },
	POSTGREY_TEXT: { group: 'Security & Anti-Spam', desc: 'SMTP message returned while a mail is greylisted.', default: 'Delayed by Postgrey' },
	POSTFIX_REJECT_UNKNOWN_CLIENT_HOSTNAME: { group: 'Security & Anti-Spam', desc: 'Reject senders whose client IP has no valid PTR/forward DNS.', default: '0' },
	SPOOF_PROTECTION: { group: 'Security & Anti-Spam', desc: 'Stop users sending as addresses they do not own.', default: '0' },
	ENABLE_SRS: { group: 'Security & Anti-Spam', desc: 'Sender Rewriting Scheme (for forwarding).', default: '0' },
	SRS_SENDER_CLASSES: { group: 'Security & Anti-Spam', desc: 'Which sender addresses SRS rewrites (envelope_sender / header_sender / both).', default: 'envelope_sender' },
	SRS_EXCLUDE_DOMAINS: { group: 'Security & Anti-Spam', desc: 'Comma-separated domains excluded from SRS rewriting.' },
	SRS_SECRET: { group: 'Security & Anti-Spam', desc: 'Comma-separated secret(s) for SRS address signing (auto-generated if unset).' },
	ENABLE_MTA_STS: { group: 'Security & Anti-Spam', desc: 'MTA-STS for outbound transport security.', default: '0' },
	POSTSCREEN_ACTION: { group: 'Security & Anti-Spam', desc: 'What Postscreen does on a violation (enforce / drop / ignore).', default: 'enforce' },
	SPAM_SUBJECT: { group: 'Security & Anti-Spam', desc: 'Subject prefix added to messages detected as spam (empty = no prefix).' },
	MOVE_SPAM_TO_JUNK: { group: 'Security & Anti-Spam', desc: 'Deliver detected spam to Junk instead of Inbox.', default: '1' },
	MARK_SPAM_AS_READ: { group: 'Security & Anti-Spam', desc: 'Mark detected spam as already read.', default: '0' },

	// Rspamd
	ENABLE_RSPAMD: { group: 'Rspamd', desc: 'Use Rspamd as the spam/malware filter.', default: '0' },
	ENABLE_RSPAMD_REDIS: { group: 'Rspamd', desc: 'Run the bundled Redis for Rspamd.', default: 'matches ENABLE_RSPAMD' },
	RSPAMD_CHECK_AUTHENTICATED: { group: 'Rspamd', desc: 'Also scan mail from authenticated senders.', default: '0' },
	RSPAMD_GREYLISTING: { group: 'Rspamd', desc: 'Rspamd greylisting module.', default: '0' },
	RSPAMD_LEARN: { group: 'Rspamd', desc: 'Autolearn ham/spam via Junk-folder Sieve.', default: '0' },
	RSPAMD_HFILTER: { group: 'Rspamd', desc: 'Hostname-validation (Hfilter) module.', default: '1' },
	RSPAMD_HFILTER_HOSTNAME_UNKNOWN_SCORE: { group: 'Rspamd', desc: 'Score added for the Hfilter "unknown hostname" symbol.', default: '6' },
	RSPAMD_NEURAL: { group: 'Rspamd', desc: 'Experimental neural-network filtering.', default: '0' },

	// SpamAssassin
	ENABLE_SPAMASSASSIN: { group: 'SpamAssassin', desc: 'SpamAssassin spam detection.', default: '0' },
	ENABLE_SPAMASSASSIN_KAM: { group: 'SpamAssassin', desc: 'Add the KAM third-party ruleset.', default: '0' },
	SPAMASSASSIN_SPAM_TO_INBOX: { group: 'SpamAssassin', desc: 'Deliver spam (tagged) instead of bouncing.', default: '1' },
	SA_TAG: { group: 'SpamAssassin', desc: 'Score at which spam info headers are added.', default: '2.0' },
	SA_TAG2: { group: 'SpamAssassin', desc: 'Score at which a message is flagged spam.', default: '6.31' },
	SA_KILL: { group: 'SpamAssassin', desc: 'Score that triggers quarantine/kill.', default: '10.0' },

	// Mail Fetching
	ENABLE_FETCHMAIL: { group: 'Mail Fetching', desc: 'Pull mail from external accounts via fetchmail.', default: '0' },
	FETCHMAIL_POLL: { group: 'Mail Fetching', desc: 'Fetchmail poll interval (seconds).', default: '300' },
	FETCHMAIL_PARALLEL: { group: 'Mail Fetching', desc: 'Run parallel fetchmail instances.', default: '0' },
	ENABLE_GETMAIL: { group: 'Mail Fetching', desc: 'Pull external mail via getmail.', default: '0' },
	GETMAIL_POLL: { group: 'Mail Fetching', desc: 'Getmail poll interval (minutes).', default: '5' },

	// Reporting & Logs
	PFLOGSUMM_TRIGGER: { group: 'Reporting & Logs', desc: 'When to send Postfix log summaries (daily_cron / logrotate).' },
	PFLOGSUMM_RECIPIENT: { group: 'Reporting & Logs', desc: 'Recipient of pflogsumm reports.' },
	PFLOGSUMM_SENDER: { group: 'Reporting & Logs', desc: 'Sender address for pflogsumm reports (falls back to REPORT_SENDER).' },
	LOGWATCH_INTERVAL: { group: 'Reporting & Logs', desc: 'Logwatch report schedule (none / daily / weekly).', default: 'none' },
	LOGWATCH_RECIPIENT: { group: 'Reporting & Logs', desc: 'Recipient of logwatch reports.' },
	LOGWATCH_SENDER: { group: 'Reporting & Logs', desc: 'Sender address for logwatch reports (falls back to REPORT_SENDER).' },
	REPORT_RECIPIENT: { group: 'Reporting & Logs', desc: 'Default recipient for all reports.' },
	REPORT_SENDER: { group: 'Reporting & Logs', desc: 'Default sender for all reports.' },
	LOGROTATE_INTERVAL: { group: 'Reporting & Logs', desc: 'Log rotation frequency (weekly / daily / monthly).', default: 'weekly' },
	LOGROTATE_COUNT: { group: 'Reporting & Logs', desc: 'Rotated logs to retain.', default: '4' },
	VIRUSMAILS_DELETE_DELAY: { group: 'Reporting & Logs', desc: 'Days to keep quarantined virus mail.', default: '7' },

	// Relay
	DEFAULT_RELAY_HOST: { group: 'Relay', desc: 'Single relay host for all outbound mail (host:port).' },
	RELAY_HOST: { group: 'Relay', desc: 'Default relay host for per-domain relaying (supports opt-out).' },
	RELAY_PORT: { group: 'Relay', desc: 'Port of the relay host.', default: '25' },
	RELAY_USER: { group: 'Relay', desc: 'Username for relay authentication.' },
	RELAY_PASSWORD: { group: 'Relay', desc: 'Password for relay authentication.' },

	// TLS / SSL
	SSL_TYPE: { group: 'TLS / SSL', desc: 'Certificate source (letsencrypt / manual / custom / self-signed).' },
	TLS_LEVEL: { group: 'TLS / SSL', desc: 'Cipher-suite strictness (modern / intermediate).', default: 'modern' },
	SSL_CERT_PATH: { group: 'TLS / SSL', desc: 'Path to the certificate file (SSL_TYPE=manual).' },
	SSL_KEY_PATH: { group: 'TLS / SSL', desc: 'Path to the private key file (SSL_TYPE=manual).' },
	SSL_ALT_CERT_PATH: { group: 'TLS / SSL', desc: 'Path to an alternate certificate for dual-cert (e.g. RSA + ECDSA).' },
	SSL_ALT_KEY_PATH: { group: 'TLS / SSL', desc: 'Path to the alternate certificate private key.' },

	// LDAP & External Auth
	ENABLE_SASLAUTHD: { group: 'LDAP & External Auth', desc: 'saslauthd authentication daemon.', default: '0' },
	SASLAUTHD_MECHANISMS: { group: 'LDAP & External Auth', desc: 'saslauthd auth mechanism (ldap / rimap / shadow / pam).' },
	SASLAUTHD_MECH_OPTIONS: { group: 'LDAP & External Auth', desc: 'Options passed to the chosen saslauthd mechanism.' },
	SASLAUTHD_LDAP_SERVER: { group: 'LDAP & External Auth', desc: 'LDAP server URI for saslauthd.' },
	SASLAUTHD_LDAP_BIND_DN: { group: 'LDAP & External Auth', desc: 'Bind DN for saslauthd LDAP.' },
	SASLAUTHD_LDAP_PASSWORD: { group: 'LDAP & External Auth', desc: 'Bind password for saslauthd LDAP.' },
	SASLAUTHD_LDAP_SEARCH_BASE: { group: 'LDAP & External Auth', desc: 'Search base for saslauthd LDAP queries.' },
	SASLAUTHD_LDAP_FILTER: { group: 'LDAP & External Auth', desc: 'LDAP filter for saslauthd user lookups.' },
	SASLAUTHD_LDAP_START_TLS: { group: 'LDAP & External Auth', desc: 'Use STARTTLS for the saslauthd LDAP connection.' },
	SASLAUTHD_LDAP_TLS_CHECK_PEER: { group: 'LDAP & External Auth', desc: 'Verify the saslauthd LDAP server certificate.' },
	SASLAUTHD_LDAP_TLS_CACERT_FILE: { group: 'LDAP & External Auth', desc: 'CA certificate file for saslauthd LDAP TLS.' },
	SASLAUTHD_LDAP_TLS_CACERT_DIR: { group: 'LDAP & External Auth', desc: 'CA certificate directory for saslauthd LDAP TLS.' },
	SASLAUTHD_LDAP_PASSWORD_ATTR: { group: 'LDAP & External Auth', desc: 'LDAP attribute holding the password.' },
	SASLAUTHD_LDAP_AUTH_METHOD: { group: 'LDAP & External Auth', desc: 'saslauthd LDAP auth method (bind / fastbind / custom).' },
	SASLAUTHD_LDAP_MECH: { group: 'LDAP & External Auth', desc: 'SASL mechanism used for the saslauthd LDAP bind.' },
	ENABLE_OAUTH2: { group: 'LDAP & External Auth', desc: 'OAuth2 (XOAUTH2/OAUTHBEARER) authentication support.' },
	OAUTH2_INTROSPECTION_URL: { group: 'LDAP & External Auth', desc: 'OAuth2 userinfo/introspection endpoint.' },
	LDAP_START_TLS: { group: 'LDAP & External Auth', desc: 'Use STARTTLS for Postfix/Dovecot LDAP connections.' },
	LDAP_SERVER_HOST: { group: 'LDAP & External Auth', desc: 'LDAP server URI.' },
	LDAP_SEARCH_BASE: { group: 'LDAP & External Auth', desc: 'LDAP search base DN.' },
	LDAP_BIND_DN: { group: 'LDAP & External Auth', desc: 'LDAP bind account DN.' },
	LDAP_BIND_PW: { group: 'LDAP & External Auth', desc: 'LDAP bind account password.' },
	LDAP_QUERY_FILTER_USER: { group: 'LDAP & External Auth', desc: 'LDAP filter for user account lookups.' },
	LDAP_QUERY_FILTER_GROUP: { group: 'LDAP & External Auth', desc: 'LDAP filter for group lookups.' },
	LDAP_QUERY_FILTER_ALIAS: { group: 'LDAP & External Auth', desc: 'LDAP filter for alias lookups.' },
	LDAP_QUERY_FILTER_DOMAIN: { group: 'LDAP & External Auth', desc: 'LDAP filter for mail-domain lookups.' },
	DOVECOT_TLS: { group: 'LDAP & External Auth', desc: 'Use TLS for Dovecot↔LDAP connections.' },
	DOVECOT_USER_FILTER: { group: 'LDAP & External Auth', desc: 'Dovecot LDAP filter for locating user objects.' },
	DOVECOT_PASS_FILTER: { group: 'LDAP & External Auth', desc: 'Dovecot LDAP filter for password verification.' },
	DOVECOT_AUTH_BIND: { group: 'LDAP & External Auth', desc: 'Use LDAP bind authentication in Dovecot.' }
};

// Supervisord services. `env` is the variable that governs the optional ones;
// `core` services should always be running. Every optional DMS service is
// listed even if it is not enabled on this host, so the dashboard can still
// describe it should it ever appear.
export const SERVICE_META = {
	postfix: { core: true, desc: 'Postfix MTA: sends and receives SMTP mail. Core service.' },
	dovecot: { core: true, desc: 'Dovecot: IMAP/POP3 access and local (LMTP) delivery. Core service.' },
	rspamd: { env: 'ENABLE_RSPAMD', desc: 'Rspamd spam/malware filter and DKIM signing.' },
	'rspamd-redis': { env: 'ENABLE_RSPAMD_REDIS', desc: 'Redis backing store for Rspamd statistics/Bayes.' },
	redis: { env: 'ENABLE_RSPAMD_REDIS', desc: 'Redis backing store for Rspamd.' },
	clamav: { env: 'ENABLE_CLAMAV', desc: 'ClamAV antivirus scanner. Heavy (~1 GB RAM); leave off unless needed.' },
	amavis: { env: 'ENABLE_AMAVIS', desc: 'Amavis content-filter front-end gluing SpamAssassin/ClamAV to Postfix.' },
	spamassassin: { env: 'ENABLE_SPAMASSASSIN', desc: 'SpamAssassin spam classifier.' },
	opendkim: { env: 'ENABLE_OPENDKIM', desc: 'OpenDKIM signing. Off automatically when Rspamd handles DKIM.' },
	opendmarc: { env: 'ENABLE_OPENDMARC', desc: 'OpenDMARC policy verification.' },
	'policyd-spf': { env: 'ENABLE_POLICYD_SPF', desc: 'SPF policy checks for inbound mail.' },
	fail2ban: { env: 'ENABLE_FAIL2BAN', desc: 'Bans IPs after repeated authentication failures.' },
	fetchmail: { env: 'ENABLE_FETCHMAIL', desc: 'Pulls mail from external accounts.' },
	getmail: { env: 'ENABLE_GETMAIL', desc: 'Alternative external-mail retrieval.' },
	postgrey: { env: 'ENABLE_POSTGREY', desc: 'Greylisting: temporarily defers unknown senders.' },
	postsrsd: { env: 'ENABLE_SRS', desc: 'Sender Rewriting Scheme daemon for forwarded mail.' },
	saslauthd: { env: 'ENABLE_SASLAUTHD', desc: 'External SASL authentication daemon.' },
	saslauthd_ldap: { env: 'ENABLE_SASLAUTHD', desc: 'saslauthd running with the LDAP mechanism (SASLAUTHD_MECHANISMS=ldap).' },
	saslauthd_rimap: { env: 'ENABLE_SASLAUTHD', desc: 'saslauthd running with the rimap mechanism (auth against a remote IMAP server).' },
	'mta-sts-daemon': { env: 'ENABLE_MTA_STS', desc: 'Serves/refreshes MTA-STS policy.' },
	'update-check': { env: 'ENABLE_UPDATE_CHECK', desc: 'Periodic check for newer DMS releases.' },
	changedetector: { core: true, desc: 'Watches config files and reloads services on change. Core service.' },
	cron: { core: true, desc: 'Scheduled maintenance (log rotation, reports, cleanup). Core service.' },
	rsyslog: { core: true, desc: 'System logging. Core service.' },
	mailserver: { core: true, desc: 'Supervisor process group for the mail stack.' }
};

// serviceContext derives a contextual status from the live service state plus
// the effective environment, so the operator does not have to cross-reference
// env vars themselves. Tones are deliberately distinct so the three states are
// never confused at a glance:
//   running:  { tone:'success', label:'running'  }  green
//   disabled: { tone:'warning', label:'disabled' }  amber  (expected, env off)
//   stopped:  { tone:'danger',  label:'stopped'  }  red    (unexpected)
//   unknown:  { tone:'muted',   label:'unknown'  }  gray   (not recognised)
export function serviceContext(svc, envMap) {
	const meta = SERVICE_META[svc.name];
	const running = svc.state === 'RUNNING';
	if (running) {
		return {
			tone: 'success',
			label: 'running',
			note: meta?.desc ?? 'Running.'
		};
	}
	if (meta?.core) {
		return {
			tone: 'danger',
			label: 'stopped',
			note: (meta.desc ? meta.desc + ' ' : '') + 'This is a core service and should be running; investigate.'
		};
	}
	if (meta?.env) {
		const enabled = truthy(envMap[meta.env]);
		if (enabled) {
			return {
				tone: 'danger',
				label: 'stopped',
				note: `${meta.desc} ${meta.env}=${envMap[meta.env]} so it should be running; investigate.`
			};
		}
		return {
			tone: 'warning',
			label: 'disabled',
			note: `${meta.desc} Disabled because ${meta.env} is not set; this is expected.`
		};
	}
	return {
		tone: 'muted',
		label: 'unknown',
		note: meta?.desc ?? 'Not a recognised docker-mailserver service; shown as-is.'
	};
}
