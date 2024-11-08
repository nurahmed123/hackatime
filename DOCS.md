<!-- omit in toc -->
# Hackatime Docs

- [Hosting](#hosting)
  - [🐳 Option 1: Use Docker](#-option-1-use-docker)
  - [🧑‍💻 Option 2: Compile and run from source](#-option-2-compile-and-run-from-source)
- [🔧 Configuration options](#-configuration-options)
  - [Supported databases](#supported-databases)
- [🔐 Authentication](#-authentication)
  - [👍 Best practices](#-best-practices)
  - [🤝 Integrations](#-integrations)
    - [Prometheus export](#prometheus-export)
      - [Scrape config example](#scrape-config-example)
      - [Grafana](#grafana)
- [🤓 Developer notes and stuff](#-developer-notes-and-stuff)
  - [Generating Swagger docs](#generating-swagger-docs)
  - [📦 Data Export](#-data-export)
  - [🧪 Tests](#-tests)
    - [Unit tests](#unit-tests)
      - [How to run](#how-to-run)
    - [API tests](#api-tests)
      - [Prerequisites (Linux only)](#prerequisites-linux-only)
      - [How to run (Linux only)](#how-to-run-linux-only)
  - [Building web assets](#building-web-assets)
    - [Precompression](#precompression)


## Hosting

### 🐳 Option 1: Use Docker

```bash
# Create a persistent volume
$ docker volume create hackatime-data

$ SALT="$(cat /dev/urandom | LC_ALL=C tr -dc 'a-zA-Z0-9' | fold -w ${1:-32} | head -n 1)"

# Run the container
$ docker run -d \
  -p 3000:3000 \
  -e "WAKAPI_PASSWORD_SALT=$SALT" \
  -v hackatime-data:/data \
  --name hackatime \
  ghcr.io/hackclub/hackatime:latest
```

Alternatively, you can use Docker Compose (`docker compose up -d`) for a more straightforward deployment.
See [compose.yml](https://github.com/kcoderhtml/hackatimee/blob/master/compose.yml) for configuration details. If you prefer to
persist data in a local directory while using SQLite as the database, make sure to set the correct `user` option in the
Docker Compose configuration to avoid permission issues.

**Note:** By default, SQLite is used as a database. To run Hackatime in Docker with MySQL or Postgres,
see [Dockerfile](https://github.com/kcoderhtml/hackatime/blob/master/Dockerfile)
and [config.default.yml](https://github.com/kcoderhtml/hackatime/blob/master/config.default.yml) for further options.

If you want to run Hackatime on **Kubernetes**, there
is [wakapi-helm-chart](https://github.com/andreymaznyak/wakapi-helm-chart) for quick and easy deployment.

### 🧑‍💻 Option 2: Compile and run from source

```bash
# Build and install
# Alternatively: go build -o wakapi
$ go install github.com/kcoderhtml/hackatime@latest

# Get default config and customize
$ curl -o hackatim.yml https://raw.githubusercontent.com/kcoderhtml/hackatime/master/config.default.yml
$ vi Hackatim.yml

# Run it
$ ./wakapi -config hackatim.yml
```

**Note:** Check the comments in `config.yml` for best practices regarding security configuration and more.

💡 When running Hackatim standalone (without Docker), it is recommended to run it as
a [SystemD service](etc/hackatime.service).

## 🔧 Configuration options

You can specify configuration options either via a config file (default: `config.yml`, customizable through the `-c`
argument) or via environment variables. Here is an overview of all options.

| YAML key / Env. variable                                                     | Default                                          | Description                                                                                                                                                                             |
| ---------------------------------------------------------------------------- | ------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `env` /<br>`ENVIRONMENT`                                                     | `dev`                                            | Whether to use development- or production settings                                                                                                                                      |
| `app.leaderboard_enabled` /<br>`WAKAPI_LEADERBOARD_ENABLED`                  | `true`                                           | Whether to enable the public leaderboard                                                                                                                                                |
| `app.ignore_user_leaderboard_preference` /<br>`WAKAPI_IGNORE_USER_LEADERBOARD_PREFERENCE` | `false`                                          | Whether to ignore user leaderboard preferences                                                                                                                                          |
| `app.leaderboard_scope` /<br>`WAKAPI_LEADERBOARD_SCOPE`                      | `7_days`                                         | Aggregation interval for public leaderboard (see [here](https://github.com/kcoderhtml/hackatime/blob/7d156cd3edeb93af2997bd95f12933b0aabef0c9/config/config.go#L71) for allowed values) |
| `app.leaderboard_generation_time` /<br>`WAKAPI_LEADERBOARD_GENERATION_TIME`  | `0 0 6 * * *,0 0 18 * * *`                       | One or multiple times of day at which to re-calculate the leaderboard                                                                                                                   |
| `app.aggregation_time` /<br>`WAKAPI_AGGREGATION_TIME`                        | `0 15 2 * * *`                                   | Time of day at which to periodically run summary generation for all users                                                                                                               |
| `app.report_time_weekly` /<br>`WAKAPI_REPORT_TIME_WEEKLY`                    | `0 0 18 * * 5`                                   | Week day and time at which to send e-mail reports                                                                                                                                       |
| `app.data_cleanup_time` /<br>`WAKAPI_DATA_CLEANUP_TIME`                      | `0 0 6 * * 0`                                    | When to perform data cleanup operations (see `app.data_retention_months`)                                                                                                               |
| `app.import_enabled` /<br>`WAKAPI_IMPORT_ENABLED`                            | `true`                                           | Whether data imports from WakaTime or other Hackatime instances are permitted                                                                                                           |
| `app.import_batch_size` /<br>`WAKAPI_IMPORT_BATCH_SIZE`                      | `50`                                             | Size of batches of heartbeats to insert to the database during importing from external services                                                                                         |
| `app.import_backoff_min` /<br>`WAKAPI_IMPORT_BACKOFF_MIN`                    | `5`                                              | "Cooldown" period in minutes before user may attempt another data import                                                                                                                |
| `app.import_max_rate` /<br>`WAKAPI_IMPORT_MAX_RATE`                          | `24`                                             | Minimum number of hours to wait after a successful data import before user may attempt another one                                                                                      |
| `app.inactive_days` /<br>`WAKAPI_INACTIVE_DAYS`                              | `7`                                              | Number of days after which to consider a user inactive (only for metrics)                                                                                                               |
| `app.heartbeat_max_age /`<br>`WAKAPI_HEARTBEAT_MAX_AGE`                      | `4320h`                                          | Maximum acceptable age of a heartbeat (see [`ParseDuration`](https://pkg.go.dev/time#ParseDuration))                                                                                    |
| `app.custom_languages`                                                       | -                                                | Map from file endings to language names                                                                                                                                                 |
| `app.avatar_url_template` /<br>`WAKAPI_AVATAR_URL_TEMPLATE`                  | (see [`config.default.yml`](config.default.yml)) | URL template for external user avatar images (e.g. from [Dicebear](https://dicebear.com) or [Gravatar](https://gravatar.com))                                                           |
| `app.date_format` /<br>`WAKAPI_DATE_FORMAT`                                  | `Mon, 02 Jan 2006`                               | Go time format strings to format human-readable date (see [`Time.Format`](https://pkg.go.dev/time#Time.Format))                                                                         |
| `app.datetime_format` /<br>`WAKAPI_DATETIME_FORMAT`                          | `Mon, 02 Jan 2006 15:04`                         | Go time format strings to format human-readable datetime (see [`Time.Format`](https://pkg.go.dev/time#Time.Format))                                                                     |
| `app.support_contact` /<br>`WAKAPI_SUPPORT_CONTACT`                          | `hostmaster@waka.hackclub.com`                   | E-Mail address to display as a support contact on the page                                                                                                                              |
| `app.data_retention_months` /<br>`WAKAPI_DATA_RETENTION_MONTHS`              | `-1`                                             | Maximum retention period in months for user data (heartbeats) (-1 for unlimited)                                                                                                        |
| `app.max_inactive_months` /<br>`WAKAPI_MAX_INACTIVE_MONTHS`                  | `12`                                             | Maximum number of inactive months after which to delete user accounts without data (-1 for unlimited)                                                                                   |
| `server.port` /<br> `WAKAPI_PORT`                                            | `3000`                                           | Port to listen on                                                                                                                                                                       |
| `server.listen_ipv4` /<br> `WAKAPI_LISTEN_IPV4`                              | `127.0.0.1`                                      | IPv4 network address to listen on (set to `'-'` to disable IPv4)                                                                                                                        |
| `server.listen_ipv6` /<br> `WAKAPI_LISTEN_IPV6`                              | `::1`                                            | IPv6 network address to listen on (set to `'-'` to disable IPv6)                                                                                                                        |
| `server.listen_socket` /<br> `WAKAPI_LISTEN_SOCKET`                          | -                                                | UNIX socket to listen on (set to `'-'` to disable UNIX socket)                                                                                                                          |
| `server.listen_socket_mode` /<br> `WAKAPI_LISTEN_SOCKET_MODE`                | `0666`                                           | Permission mode to create UNIX socket with                                                                                                                                              |
| `server.timeout_sec` /<br> `WAKAPI_TIMEOUT_SEC`                              | `30`                                             | Request timeout in seconds                                                                                                                                                              |
| `server.tls_cert_path` /<br> `WAKAPI_TLS_CERT_PATH`                          | -                                                | Path of SSL server certificate (leave blank to not use HTTPS)                                                                                                                           |
| `server.tls_key_path` /<br> `WAKAPI_TLS_KEY_PATH`                            | -                                                | Path of SSL server private key (leave blank to not use HTTPS)                                                                                                                           |
| `server.base_path` /<br> `WAKAPI_BASE_PATH`                                  | `/`                                              | Web base path (change when running behind a proxy under a sub-path)                                                                                                                     |
| `server.public_url` /<br> `WAKAPI_PUBLIC_URL`                                | `http://localhost:3000`                          | URL at which your Hackatime instance can be found publicly                                                                                                                              |
| `security.password_salt` /<br> `WAKAPI_PASSWORD_SALT`                        | -                                                | Pepper to use for password hashing                                                                                                                                                      |
| `security.insecure_cookies` /<br> `WAKAPI_INSECURE_COOKIES`                  | `false`                                          | Whether or not to allow cookies over HTTP                                                                                                                                               |
| `security.cookie_max_age` /<br> `WAKAPI_COOKIE_MAX_AGE`                      | `172800`                                         | Lifetime of authentication cookies in seconds or `0` to use [Session](https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies#Define_the_lifetime_of_a_cookie) cookies                |
| `security.allow_signup` /<br> `WAKAPI_ALLOW_SIGNUP`                          | `true`                                           | Whether to enable user registration                                                                                                                                                     |
| `security.signup_captcha` /<br> `WAKAPI_SIGNUP_CAPTCHA`                      | `false`                                          | Whether the registration form requires solving a CAPTCHA                                                                                                                                |
| `security.invite_codes` /<br> `WAKAPI_INVITE_CODES`                          | `true`                                           | Whether to enable registration by invite codes. Primarily useful if registration is disabled (invite-only server).                                                                      |
| `security.disable_frontpage` /<br> `WAKAPI_DISABLE_FRONTPAGE`                | `false`                                          | Whether to disable landing page (useful for personal instances)                                                                                                                         |
| `security.expose_metrics` /<br> `WAKAPI_EXPOSE_METRICS`                      | `false`                                          | Whether to expose Prometheus metrics under `/api/metrics`                                                                                                                               |
| `security.trusted_header_auth` /<br> `WAKAPI_TRUSTED_HEADER_AUTH`            | `false`                                          | Whether to enable trusted header authentication for reverse proxies (see [#534](https://github.com/muety/wakatime/issues/534)). **Use with caution!**                                   |
| `security.trusted_header_auth_key` /<br> `WAKAPI_TRUSTED_HEADER_AUTH_KEY`    | `Remote-User`                                    | Header field for trusted header authentication. **Caution:** proxy must be configured to strip this header from client requests!                                                        |
| `security.trust_reverse_proxy_ips` /<br> `WAKAPI_TRUST_REVERSE_PROXY_IPS`    | -                                                | Comma-separated list of IPv4 or IPv6 addresses or CIDRs of reverse proxies to trust to handle authentication (e.g. `172.17.0.1`, `192.168.0.0/24`, `[::1]`).                            |
| `security.signup_max_rate` /<br> `WAKAPI_SIGNUP_MAX_RATE`                    | `5/1h`                                           | Rate limiting config for signup endpoint in format `<max_req>/<multiplier><unit>`, where `unit` is one of `s`, `m` or `h`.                                                              |
| `security.login_max_rate` /<br> `WAKAPI_LOGIN_MAX_RATE`                      | `10/1m`                                          | Rate limiting config for login endpoint in format `<max_req>/<multiplier><unit>`, where `unit` is one of `s`, `m` or `h`.                                                               |
| `security.password_reset_max_rate` /<br> `WAKAPI_PASSWORD_RESET_MAX_RATE`    | `5/1h`                                           | Rate limiting config for password reset endpoint in format `<max_req>/<multiplier><unit>`, where `unit` is one of `s`, `m` or `h`.                                                      |
| `db.host` /<br> `WAKAPI_DB_HOST`                                             | -                                                | Database host                                                                                                                                                                           |
| `db.port` /<br> `WAKAPI_DB_PORT`                                             | -                                                | Database port                                                                                                                                                                           |
| `db.socket` /<br> `WAKAPI_DB_SOCKET`                                         | -                                                | Database UNIX socket (alternative to `host`) (for MySQL only)                                                                                                                           |
| `db.user` /<br> `WAKAPI_DB_USER`                                             | -                                                | Database user                                                                                                                                                                           |
| `db.password` /<br> `WAKAPI_DB_PASSWORD`                                     | -                                                | Database password                                                                                                                                                                       |
| `db.name` /<br> `WAKAPI_DB_NAME`                                             | `wakapi_db.db`                                   | Database name                                                                                                                                                                           |
| `db.dialect` /<br> `WAKAPI_DB_TYPE`                                          | `sqlite3`                                        | Database type (one of `sqlite3`, `mysql`, `postgres`, `cockroach`, `mssql`)                                                                                                             |
| `db.charset` /<br> `WAKAPI_DB_CHARSET`                                       | `utf8mb4`                                        | Database connection charset (for MySQL only)                                                                                                                                            |
| `db.max_conn` /<br> `WAKAPI_DB_MAX_CONNECTIONS`                              | `2`                                              | Maximum number of database connections                                                                                                                                                  |
| `db.ssl` /<br> `WAKAPI_DB_SSL`                                               | `false`                                          | Whether to use TLS encryption for database connection (Postgres and CockroachDB only)                                                                                                   |
| `db.automgirate_fail_silently` /<br> `WAKAPI_DB_AUTOMIGRATE_FAIL_SILENTLY`   | `false`                                          | Whether to ignore schema auto-migration failures when starting up                                                                                                                       |
| `mail.enabled` /<br> `WAKAPI_MAIL_ENABLED`                                   | `true`                                           | Whether to allow Hackatime to send e-mail (e.g. for password resets) |
| `mail.welcome_enabled` /<br> `WAKAPI_WELCOME_ENABLED`                        | `true`                                           | Whether Hackatime should send an e-mail on user signup |
| `mail.sender` /<br> `WAKAPI_MAIL_SENDER`                                     | `Hackatime <noreply@waka.hackclub.com>`          | Default sender address for outgoing mails |
| `mail.provider` /<br> `WAKAPI_MAIL_PROVIDER`                                 | `smtp`                                           | Implementation to use for sending mails (one of [`smtp`])                                                                                                                               |
| `mail.smtp.host` /<br> `WAKAPI_MAIL_SMTP_HOST`                               | -                                                | SMTP server address for sending mail (if using `smtp` mail provider)                                                                                                                    |
| `mail.smtp.port` /<br> `WAKAPI_MAIL_SMTP_PORT`                               | -                                                | SMTP server port (usually 465)                                                                                                                                                          |
| `mail.smtp.username` /<br> `WAKAPI_MAIL_SMTP_USER`                           | -                                                | SMTP server authentication username                                                                                                                                                     |
| `mail.smtp.password` /<br> `WAKAPI_MAIL_SMTP_PASS`                           | -                                                | SMTP server authentication password                                                                                                                                                     |
| `mail.smtp.tls` /<br> `WAKAPI_MAIL_SMTP_TLS`                                 | `false`                                          | Whether the SMTP server requires TLS encryption (`false` for STARTTLS or no encryption)                                                                                                 |
| `mail.smtp.skip_verify` /<br> `WAKAPI_MAIL_SMTP_SKIP_VERIFY`                 | `false`                                          | Whether to allow invalid or self-signed certificates for TLS-encrypted SMTP                                                                                                             |
| `sentry.dsn` /<br> `WAKAPI_SENTRY_DSN`                                       | –                                                | DSN for to integrate [Sentry](https://sentry.io) for error logging and tracing (leave empty to disable)                                                                                 |
| `sentry.environment` /<br> `WAKAPI_SENTRY_ENVIRONMENT`                       | (`env`)                                          | Sentry [environment](https://docs.sentry.io/concepts/key-terms/environments/) tag (defaults to `env` / `ENV`)                                                                           |
| `sentry.enable_tracing` /<br> `WAKAPI_SENTRY_TRACING`                        | `false`                                          | Whether to enable Sentry request tracing                                                                                                                                                |
| `sentry.sample_rate` /<br> `WAKAPI_SENTRY_SAMPLE_RATE`                       | `0.75`                                           | Probability of tracing a request in Sentry                                                                                                                                              |
| `sentry.sample_rate_heartbeats` /<br> `WAKAPI_SENTRY_SAMPLE_RATE_HEARTBEATS` | `0.1`                                            | Probability of tracing a heartbeat request in Sentry                                                                                                                                    |
| `quick_start` /<br> `WAKAPI_QUICK_START`                                     | `false`                                          | Whether to skip initial boot tasks. Use only for development purposes!                                                                                                                  |
| `enable_pprof` /<br> `WAKAPI_ENABLE_PPROF`                                   | `false`                                          | Whether to expose [pprof](https://pkg.go.dev/runtime/pprof) profiling data as an endpoint for debugging                                                                                 |

### Supported databases

Hackatime uses [GORM](https://gorm.io) as an ORM. As a consequence, a set of different relational databases is supported.

-   [SQLite](https://sqlite.org/) (_default, easy setup_)
-   [MySQL](https://hub.docker.com/_/mysql) (_recommended, because most extensively tested_)
-   [MariaDB](https://hub.docker.com/_/mariadb) (_open-source MySQL alternative_)
-   [Postgres](https://hub.docker.com/_/postgres) (_open-source as well_)
-   [CockroachDB](https://www.cockroachlabs.com/docs/stable/install-cockroachdb-linux.html) (_cloud-native, distributed,
    Postgres-compatible API_)
-   [Microsoft SQL Server](https://hub.docker.com/_/microsoft-mssql-server) (_Microsoft SQL Server_)

## 🔐 Authentication

Hackatime supports different types of user authentication.

-   **Cookie:** This method is used in the browser. Users authenticate by sending along an encrypted, secure, HTTP-only
    cookie (`wakapi_auth`) that was set in the server's response upon login.
-   **API key:**
    -   **Via header:** This method is inspired
        by [WakaTime's auth. mechanism](https://wakatime.com/developers/#authentication) and is the common way to
        authenticate against API endpoints. Users set the `Authorization` header to `Basic <BASE64_TOKEN>`, where the
        latter part corresponds to your base64-hashed API key.
    -   **Vis query param:** Alternatively, users can also pass their plain API key as a query parameter (
        e.g. `?api_key=86648d74-19c5-452b-ba01-fb3ec70d4c2f`) in the URL with every request.
-   **Trusted header:** This mechanism allows to delegate authentication to a **reverse proxy** (e.g. for SSO), that
    Hackatime will then trust blindly. See [#534](https://github.com/muety/wakapi/issues/534) for details.
    -   Must be enabled via `trusted_header_auth` and configuring `trust_reverse_proxy_ip` in the config
    -   Warning: This type of authentication is quite prone to misconfiguration. Make sure that your reverse proxy
        properly strips relevant headers from client requests.

### 👍 Best practices

It is recommended to use wakapi behind a **reverse proxy**, like [Caddy](https://caddyserver.com)
or [nginx](https://www.nginx.com/), to enable **TLS encryption** (HTTPS).

However, if you want to expose your wakapi instance to the public anyway, you need to set `server.listen_ipv4`
to `0.0.0.0` in `config.yml`.

### 🤝 Integrations

#### Prometheus export

You can export your Hackatime statistics to Prometheus to view them in a Grafana dashboard or so. Here is how.

```bash
# 1. Start Hackatime with the feature enabled
$ export WAKAPI_EXPOSE_METRICS=true
$ ./wakapi

# 2. Get your API key and hash it
$ echo "<YOUR_API_KEY>" | base64

# 3. Add a Prometheus scrape config to your prometheus.yml (see below)
```

##### Scrape config example

```yml
# prometheus.yml
# (assuming your Hackatime instance listens at localhost, port 3000)

scrape_configs:
    - job_name: 'wakapi'
      scrape_interval: 1m
      metrics_path: '/api/metrics'
      bearer_token: '<YOUR_BASE64_HASHED_TOKEN>'
      static_configs:
          - targets: ['localhost:3000']
```

##### Grafana

There is also a [nice Grafana dashboard](https://grafana.com/grafana/dashboards/12790), provided by the author
of [wakatime_exporter](https://github.com/MacroPower/wakatime_exporter).

![](https://grafana.com/api/dashboards/12790/images/8741/image)

## 🤓 Developer notes and stuff

### Generating Swagger docs

```bash
$ go install github.com/swaggo/swag/cmd/swag@latest
$ swag init -o static/docs
```

### 📦 Data Export

You can export your coding activity from Hackatime to CSV in the form of raw heartbeats. While there is no way to
accomplish this directly through the web UI, we provide an easy-to-use Python [script](scripts/download_heartbeats.py)
instead.

```bash
$ pip install requests tqdm
$ python scripts/download_heartbeats.py --api_key API_KEY [--url URL] [--from FROM] [--to TO] [--output OUTPUT]
```

<details>

<summary>Example</summary>

```bash
python scripts/download_heartbeats.py --api_key 04648d14-15c9-432b-b901-dbeec70d4eaf \
  --url https://waka.hackclub.com/api \
  --from 2023-01-01 \
  --to 2023-01-31 \
  --output wakapi_export.csv
```

</details>

### 🧪 Tests

#### Unit tests

Unit tests are supposed to test business logic on a fine-grained level. They are implemented as part of the application,
using Go's [testing](https://pkg.go.dev/testing?utm_source=godoc) package
alongside [stretchr/testify](https://pkg.go.dev/github.com/stretchr/testify).

##### How to run

```bash
$ CGO_ENABLED=0 go test `go list ./... | grep -v 'github.com/kcoderhtml/hackatime/scripts'` -json -coverprofile=coverage/coverage.out ./... -run ./...
```

#### API tests

API tests are implemented as black box tests, which interact with a fully-fledged, standalone Hackatime through HTTP
requests. They are supposed to check Hackatime's web stack and endpoints, including response codes, headers and data on a
syntactical level, rather than checking the actual content that is returned.

Our API (or end-to-end, in some way) tests are implemented as a [Postman](https://www.postman.com/) collection and can
be run either from inside Postman, or using [newman](https://www.npmjs.com/package/newman) as a command-line runner.

To get a predictable environment, tests are run against a fresh and clean Hackatime instance with a SQLite database that is
populated with nothing but some seed data (see [data.sql](testing/data.sql)). It is usually recommended for software
tests to be [safe](https://www.restapitutorial.com/lessons/idempotency.html), stateless and without side effects. In
contrary to that paradigm, our API tests strictly require a fixed execution order (which Postman assures) and their
assertions may rely on specific previous tests having succeeded.

##### Prerequisites (Linux only)

```bash
# 1. sqlite (cli)
$ sudo apt install sqlite  # Fedora: sudo dnf install sqlite

# 2. newman
$ npm install -g newman
```

##### How to run (Linux only)

```bash
$ ./testing/run_api_tests.sh
```

### Building web assets

To keep things minimal, all JS and CSS assets are included as static files and checked in to
Git. [TailwindCSS](https://tailwindcss.com/docs/installation#building-for-production)
and [Iconify](https://iconify.design/docs/icon-bundles/) require an additional build step. To only require this at the
time of development, the compiled assets are checked in to Git as well.

```bash
$ yarn
$ yarn build  # or: yarn watch
```

New icons can be added by editing the `icons` array in [scripts/bundle_icons.js](scripts/bundle_icons.js).

#### Precompression

As explained in [#284](https://github.com/muety/wakapi/issues/284), precompressed (using Brotli) versions of some of the
assets are delivered to save additional bandwidth. This was inspired by
Caddy's [`precompressed`](https://caddyserver.com/docs/caddyfile/directives/file_server)
directive. [
`gzipped.FileServer`](https://github.com/kcoderhtml/hackatime/blob/07a367ce0a97c7738ba8e255e9c72df273fd43a3/main.go#L249)
checks for every static file's `.br` or `.gz` equivalents and, if present, delivers those instead of the actual file,
alongside `Content-Encoding: br`. Currently, compressed assets are simply checked in to Git. Later we might want to have
this be part of a new build step.

To pre-compress files, run this:

```bash
# Install brotli first
$ sudo apt install brotli  # or: sudo dnf install brotli

# Watch, build and compress
$ yarn watch:compress

# Alternatively: build and compress only
$ yarn build:all:compress

# Alternatively: compress only
$ yarn compress
```