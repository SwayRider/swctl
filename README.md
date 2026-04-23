# swctl

Management CLI tool for the SwayRider platform. Provides command-line access to the Auth, Health, Search, and Region services.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Commands](#commands)
  - [auth](#auth)
    - [auth whoami](#auth-whoami)
    - [auth get-user](#auth-get-user)
    - [auth check-password-strength](#auth-check-password-strength)
    - [auth create-admin](#auth-create-admin)
    - [auth create-user](#auth-create-user)
    - [auth change-password](#auth-change-password)
    - [auth create-service-client](#auth-create-service-client)
    - [auth list-service-clients](#auth-list-service-clients)
    - [auth delete-service-client](#auth-delete-service-client)
  - [health](#health)
    - [health ping](#health-ping)
    - [health check](#health-check)
  - [search](#search)
    - [search geocode](#search-geocode)
    - [search reverse](#search-reverse)
  - [region](#region)
    - [region search-point](#region-search-point)
    - [region search-box](#region-search-box)
    - [region find-region-path](#region-find-region-path)

---

## Installation

```bash
go build -o swctl ./cmd/swctl
```

Or run directly:

```bash
go run ./cmd/swctl [command]
```

---

## Configuration

`swctl` resolves configuration in the following order (highest priority first):

1. **CLI flags** — passed directly on the command line
2. **Environment variables** — set in your shell
3. **`swctl.conf`** — a `.env`-style file in the working directory

### swctl.conf

Copy and edit `swctl.conf` to set defaults for your environment:

```env
# AuthService
AUTH_HOST=127.0.0.1
AUTH_PORT=8081
AUTH_USER=admin@example.com
AUTH_PASSWORD=yourpassword

# HealthService (target service host/port, not a dedicated health service)
HEALTH_HOST=127.0.0.1
HEALTH_PORT=

# SearchService
SEARCH_HOST=127.0.0.1
SEARCH_PORT=

# RegionService
REGION_HOST=127.0.0.1
REGION_PORT=
```

All commands that interact with the AuthService accept `--auth-host` and `--auth-port` flags (or `AUTH_HOST` / `AUTH_PORT` env vars). Commands for other services follow the same pattern with their own prefix.

---

## Commands

### auth

Interact with the AuthService.

```
swctl auth [--auth-host HOST] [--auth-port PORT] <command>
```

| Flag | Env var | Description |
|---|---|---|
| `--auth-host` | `AUTH_HOST` | Host of the AuthService (required) |
| `--auth-port` | `AUTH_PORT` | Port of the AuthService (required) |

---

#### auth whoami

Show the identity resolved by the configured credentials. Useful for verifying that `AUTH_USER`/`AUTH_PASSWORD` are correct and the service is reachable.

**Alias:** `wai`

```
swctl auth whoami [--user USER] [--password PASSWORD]
```

| Flag | Env var | Description |
|---|---|---|
| `--user`, `-u` | `AUTH_USER` | Admin email (required) |
| `--password`, `-p` | `AUTH_PASSWORD` | Admin password (required) |

**Example:**

```bash
swctl auth whoami
```

**Output:**

```
    Email: admin@example.com
    UserID: 3f2e1a...
    Account Type: admin
    Verified: true
    Admin: true
```

---

#### auth get-user

Look up a user by email address or user ID.

**Alias:** `gu`

```
swctl auth get-user <email|userId> [--user USER] [--password PASSWORD]
```

| Argument | Description |
|---|---|
| `email\|userId` | Email address (contains `@`) or UUID user ID |

| Flag | Env var | Description |
|---|---|---|
| `--user`, `-u` | `AUTH_USER` | Admin email (required) |
| `--password`, `-p` | `AUTH_PASSWORD` | Admin password (required) |

**Examples:**

```bash
swctl auth get-user alice@example.com
swctl auth get-user 3f2e1a00-dead-beef-0000-000000000001
```

---

#### auth check-password-strength

Check whether a password meets the platform's strength requirements. Does not require authentication.

**Alias:** `cps`

```
swctl auth check-password-strength <password>
```

| Argument | Description |
|---|---|
| `password` | The password to evaluate |

**Example:**

```bash
swctl auth check-password-strength "hunter2"
```

**Output:** A human-readable message from the service (e.g. `Password is too weak`).

---

#### auth create-admin

Create a new admin user. Requires existing admin credentials to authenticate.

**Alias:** `ca`

```
swctl auth create-admin <email> <password> [--user USER] [--password PASSWORD]
```

| Argument | Description |
|---|---|
| `email` | Email for the new admin account |
| `password` | Password for the new admin account |

| Flag | Env var | Description |
|---|---|---|
| `--user`, `-u` | `AUTH_USER` | Authenticating admin email (required) |
| `--password`, `-p` | `AUTH_PASSWORD` | Authenticating admin password (required) |

**Example:**

```bash
swctl auth create-admin newadmin@example.com SecurePass123!
```

---

#### auth create-user

Register a new regular user. Optionally mark them as verified and set their account type.

**Alias:** `cu`

```
swctl auth create-user <email> <password> [flags]
```

| Argument | Description |
|---|---|
| `email` | Email for the new user |
| `password` | Password for the new user |

| Flag | Env var | Default | Description |
|---|---|---|---|
| `--user`, `-u` | `AUTH_USER` | — | Authenticating admin email (required) |
| `--password`, `-p` | `AUTH_PASSWORD` | — | Authenticating admin password (required) |
| `--verified`, `-v` | — | `false` | Mark the user's email as verified immediately |
| `--account-type`, `-t` | — | `free` | Account type (e.g. `free`, `premium`) |

**Examples:**

```bash
# Create a basic free user
swctl auth create-user user@example.com MyPass123!

# Create a verified premium user
swctl auth create-user user@example.com MyPass123! --verified --account-type premium
```

---

#### auth change-password

Change the password for the authenticated user. The `--password` flag provides the **current** (old) password; the new password is passed as an argument.

**Alias:** `chp`

```
swctl auth change-password <newPassword> [--user USER] [--password CURRENT_PASSWORD]
```

| Argument | Description |
|---|---|
| `newPassword` | The new password to set |

| Flag | Env var | Description |
|---|---|---|
| `--user`, `-u` | `AUTH_USER` | User email (required) |
| `--password`, `-p` | `AUTH_PASSWORD` | Current password — used for login and as the old password (required) |

**Example:**

```bash
swctl auth change-password NewSecurePass456! --user me@example.com --password OldPass123!
```

---

#### auth create-service-client

Create an OAuth2-style service client with a set of scopes. The client secret is displayed once — store it immediately.

**Alias:** `csc`

```
swctl auth create-service-client <name> <scope...> [flags]
```

| Argument | Description |
|---|---|
| `name` | Name of the service client |
| `scope...` | One or more scopes to grant (space-separated) |

| Flag | Env var | Description |
|---|---|---|
| `--user`, `-u` | `AUTH_USER` | Admin email (required) |
| `--password`, `-p` | `AUTH_PASSWORD` | Admin password (required) |
| `--description`, `-d` | — | Human-readable description of the client |

**Example:**

```bash
swctl auth create-service-client my-service read:tiles write:routes \
    --description "Routing microservice"
```

**Output:**

```
WARNING: Store the client secret securely — it will not be shown again.
client id: abc123...
client secret: xyz789...
```

---

#### auth list-service-clients

List all registered service clients in a formatted table.

**Alias:** `lsc`

```
swctl auth list-service-clients [flags]
```

| Flag | Env var | Default | Description |
|---|---|---|---|
| `--user`, `-u` | `AUTH_USER` | — | Admin email (required) |
| `--password`, `-p` | `AUTH_PASSWORD` | — | Admin password (required) |
| `--page` | — | `0` | Page number (`0` = all) |
| `--page-size` | — | `0` | Results per page (`0` = all) |

**Example:**

```bash
swctl auth list-service-clients
swctl auth list-service-clients --page 1 --page-size 20
```

**Output:**

```
+--------------------+----------------------------------+----------------+------------------------------------------------------------------+
| NAME               | CLIENTID                         | SCOPES         | DESCRIPTION                                                      |
+--------------------+----------------------------------+----------------+------------------------------------------------------------------+
| my-service         | abc123...                        | read:tiles     | Routing microservice                                             |
+--------------------+----------------------------------+----------------+------------------------------------------------------------------+
```

---

#### auth delete-service-client

Delete a service client by its client ID.

**Alias:** `dsc`

```
swctl auth delete-service-client <clientId> [--user USER] [--password PASSWORD]
```

| Argument | Description |
|---|---|
| `clientId` | The client ID of the service client to delete |

| Flag | Env var | Description |
|---|---|---|
| `--user`, `-u` | `AUTH_USER` | Admin email (required) |
| `--password`, `-p` | `AUTH_PASSWORD` | Admin password (required) |

**Example:**

```bash
swctl auth delete-service-client abc123def456...
```

---

### health

Check the availability and health status of any SwayRider service. Point it at any service's host and port — all services expose a health endpoint.

```
swctl health [--health-host HOST] [--health-port PORT] <command>
```

| Flag | Env var | Description |
|---|---|---|
| `--health-host` | `HEALTH_HOST` | Host of the target service (required) |
| `--health-port` | `HEALTH_PORT` | Port of the target service (required) |

---

#### health ping

Ping a service to confirm it is reachable. Exits with a non-zero status if the service does not respond.

```
swctl health ping
```

**Example:**

```bash
# Check if the AuthService is up
swctl health ping --health-host 127.0.0.1 --health-port 8081
```

**Output:** `pong`

---

#### health check

Check the health status of a service, optionally scoped to a specific internal component.

```
swctl health check [--component COMPONENT]
```

| Flag | Description |
|---|---|
| `--component`, `-c` | Internal component name to check (omit for overall service status) |

**Example:**

```bash
# Overall service health
swctl health check --health-host 127.0.0.1 --health-port 8081

# Check a specific component
swctl health check --health-host 127.0.0.1 --health-port 8081 --component database
```

**Output:** `status: UP` (possible values: `UP`, `DOWN`, `UNKNOWN`)

---

### search

Geocode text queries and reverse geocode coordinates via the SearchService. Requires authentication through the AuthService.

```
swctl search \
    [--auth-host HOST] [--auth-port PORT] \
    [--search-host HOST] [--search-port PORT] \
    <command>
```

| Flag | Env var | Description |
|---|---|---|
| `--auth-host` | `AUTH_HOST` | Host of the AuthService (required) |
| `--auth-port` | `AUTH_PORT` | Port of the AuthService (required) |
| `--search-host` | `SEARCH_HOST` | Host of the SearchService (required) |
| `--search-port` | `SEARCH_PORT` | Port of the SearchService (required) |

---

#### search geocode

Convert a text query into geographic coordinates.

**Alias:** `gc`

```
swctl search geocode <query> [flags]
```

| Argument | Description |
|---|---|
| `query` | Text to geocode (e.g. `"Brussels, Belgium"`) |

| Flag | Env var | Default | Description |
|---|---|---|---|
| `--user`, `-u` | `AUTH_USER` | — | Auth email (required) |
| `--password`, `-p` | `AUTH_PASSWORD` | — | Auth password (required) |
| `--size` | — | `10` | Maximum number of results to return |
| `--lang` | — | — | Preferred language code for result labels (e.g. `en`, `nl`, `fr`) |
| `--lat` | — | — | Focus point latitude (biases results toward this location) |
| `--lon` | — | — | Focus point longitude |

**Example:**

```bash
swctl search geocode "Grote Markt, Antwerp"
swctl search geocode "coffee" --lat 51.2194 --lon 4.4025 --size 5 --lang nl
```

**Output:**

```
    Grote Markt, Antwerp, Belgium [venue] lat=51.221117 lon=4.399708 confidence=1.00
    Grote Markt, Antwerp, Belgium [address] lat=51.221044 lon=4.399572 confidence=0.92
```

---

#### search reverse

Convert a coordinate pair into a place name.

**Alias:** `rv`

```
swctl search reverse <lat> <lon> [flags]
```

| Argument | Description |
|---|---|
| `lat` | Latitude |
| `lon` | Longitude |

| Flag | Env var | Default | Description |
|---|---|---|---|
| `--user`, `-u` | `AUTH_USER` | — | Auth email (required) |
| `--password`, `-p` | `AUTH_PASSWORD` | — | Auth password (required) |
| `--size` | — | `5` | Maximum number of results to return |
| `--lang` | — | — | Preferred language code for result labels |

**Example:**

```bash
swctl search reverse 51.2194 4.4025
swctl search reverse 50.8503 4.3517 --lang fr --size 3
```

**Output:**

```
    Antwerpen, Antwerp, Belgium [locality] lat=51.219405 lon=4.402339 confidence=1.00
```

---

### region

Query geographic regions via the RegionService. No authentication required.

```
swctl region [--region-host HOST] [--region-port PORT] <command>
```

| Flag | Env var | Description |
|---|---|---|
| `--region-host` | `REGION_HOST` | Host of the RegionService (required) |
| `--region-port` | `REGION_PORT` | Port of the RegionService (required) |

---

#### region search-point

Find which geographic regions contain a given coordinate.

**Alias:** `sp`

```
swctl region search-point <lat> <lon> [--extended]
```

| Argument | Description |
|---|---|
| `lat` | Latitude |
| `lon` | Longitude |

| Flag | Default | Description |
|---|---|---|
| `--extended`, `-e` | `false` | Include extended (non-core) region data in the output |

**Example:**

```bash
swctl region search-point 51.2194 4.4025
swctl region search-point 51.2194 4.4025 --extended
```

**Output:**

```
Core regions:     BE, BE-VAN, BE-VAN-ANT
Extended regions: EU, BENELUX
```

---

#### region search-box

Find which geographic regions intersect a bounding box.

**Alias:** `sb`

```
swctl region search-box <minLat> <minLon> <maxLat> <maxLon> [--extended]
```

| Argument | Description |
|---|---|
| `minLat` | Bottom-left latitude |
| `minLon` | Bottom-left longitude |
| `maxLat` | Top-right latitude |
| `maxLon` | Top-right longitude |

| Flag | Default | Description |
|---|---|---|
| `--extended`, `-e` | `false` | Include extended region data |

**Example:**

```bash
swctl region search-box 50.5 2.5 51.5 6.5
```

**Output:**

```
Core regions:     BE, BE-VAN, BE-WAL, BE-BRU, NL, NL-ZB, NL-NB
```

---

#### region find-region-path

Find the ordered list of regions that form a path between two region codes.

**Alias:** `frp`

```
swctl region find-region-path <fromRegion> <toRegion>
```

| Argument | Description |
|---|---|
| `fromRegion` | Starting region code |
| `toRegion` | Destination region code |

**Example:**

```bash
swctl region find-region-path BE-VAN-ANT NL-ZH-RTD
```

**Output:**

```
BE-VAN-ANT → BE-VAN → BE → NL → NL-ZH → NL-ZH-RTD
```
