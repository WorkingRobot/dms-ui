# dms-ui

[![Build](https://img.shields.io/github/actions/workflow/status/WorkingRobot/dms-ui/build-web.yml?style=for-the-badge&label=Build)](https://github.com/WorkingRobot/dms-ui/pkgs/container/dms-ui)
[![Image](https://img.shields.io/badge/ghcr.io-workingrobot%2Fdms--ui-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://github.com/WorkingRobot/dms-ui/pkgs/container/dms-ui)
[![License](https://img.shields.io/github/license/WorkingRobot/dms-ui?style=for-the-badge)](/LICENSE)

A lightweight web admin panel for [docker-mailserver](https://github.com/docker-mailserver/docker-mailserver) (DMS).

DMS has no HTTP API, so all administration goes through its in-container `setup`. `dms-ui` is a single Go binary (with an embedded SvelteKit SPA) that drives `setup` over a mounted Docker socket, so Postfix/Dovecot stay consistent and DMS's change detector reloads gracefully (no mail disruption).

> [!NOTE]
> A significant portion of this project was built by Claude in an afternoon. Don't use it if you don't trust it with your mail.

## Features

- **Dashboard:** env-aware service status, rspamd stats, account/alias/queue counts
- **Accounts:** add / delete (optional mailbox purge) / change password, per-mailbox storage & per-folder message counts, quotas, send/receive restrictions
- **Aliases:** virtual alias management
- **Dovecot masters:** SASL master users
- **DKIM:** list keys with ready-to-publish DNS TXT records, generate / rotate (rspamd)
- **Relay:** relay domains, SASL credentials, per-domain exclusions
- **Fail2ban:** per-jail status, ban / unban
- **Logs:** mail log, fail2ban log, Postfix queue
- **Environment:** the full DMS config surface, grouped & documented

## Quick Start

> [!CAUTION]
> `dms-ui` has **no built-in authentication** and is root-equivalent (it talks
> to the Docker socket). Run it behind a trusted authenticating reverse proxy
> (e.g. Traefik + Authentik); you are responsible for that exposure.

```bash
docker run -d --name dms-ui \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e MAIL_CONTAINER=mailserver \
  -p 127.0.0.1:8080:8080 \
  ghcr.io/workingrobot/dms-ui:edge
```

See [`deploy/`](/deploy) for a compose snippet and a Traefik + Authentik example.

## Configuration

| Var | Default | Notes |
|---|---|---|
| `MAIL_CONTAINER` | `mailserver` | Target docker-mailserver container name |
| `LISTEN_ADDR` | `:8080` | HTTP bind |
| `PROXY_EMAIL_HEADER` | `X-Forwarded-Email` | Header the proxy sets to identify the admin |
| `DOCKER_HOST` | _unset_ | Honoured by the docker client (e.g. `ssh://host` for remote dev) |

## Develop

```sh
make build          # builds web + binary
make run            # runs on 127.0.0.1:8099 (override MAIL_CONTAINER/DOCKER_HOST)
# or, with live reload of the UI:
cd web && npm run dev   # proxies /api to 127.0.0.1:8099
```

## Contributing

Contributions, bug reports, and feature requests are welcome! Please open an [issue](https://github.com/WorkingRobot/dms-ui/issues) or a [pull request](https://github.com/WorkingRobot/dms-ui/pulls).
