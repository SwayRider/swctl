# Code Review — 2026-08-19

Security-focused review of the full current codebase of `swctl` (not diff-based). No prior `review/CODE_REVIEW_*.md` exists in this repo, so all findings below are new. See [`Docs/REVIEW.md`](../../Docs/REVIEW.md) for how findings in this file are tracked and marked fixed.

No findings on command injection (swctl never shells out — no `os/exec` usage found) or SQL injection (no direct DB access; all operations go through gRPC service clients). Privilege scope is reasonable — it only requires whatever authservice admin scopes the operator already has.

### 1. ~~Admin credentials committed to git history~~ — NOT APPLICABLE 2026-08-19

`swctl.conf` (gitignored today) was tracked in the first two commits of the pushed `main` branch: `3cc7faa` (initial commit) and `f3d2566`. `git show 3cc7faa:swctl.conf` returns a credential pair — `AUTH_USER=maarten@hevanto-it.com`, `AUTH_PASSWORD=maarten123`. The file was later removed from tracking and gitignored (`b069bb1 Add gitignore`), but the blob remains permanently reachable via `git log --all`/`git show` on the remote.

Confirmed by the user (maarten@hevanto.be): this credential belongs to an internal test environment, not production — accepted risk, no rotation or history rewrite needed. Kept here (rather than deleted) per [`Docs/REVIEW.md`](../../Docs/REVIEW.md) so the pattern (real-looking credentials having lived in `swctl.conf` before it was gitignored) stays visible for future reviews, and so a *future* credential added to this file is understood to need the same scrutiny.

### 2. ~~Admin/user passwords passed as CLI positional arguments~~ — FIXED 2026-08-19

`auth create-admin <email> <password>`, `auth create-user <email> <password>`, and `auth change-password <newPassword>` (README.md:196,224,258; `internal/cmd/auth/create_admin.go`, `create_user.go`, `change_password.go`) take plaintext passwords as positional args, in addition to `--password`/`-p`/`AUTH_PASSWORD` for the authenticating admin. Both forms land in shell history (`~/.zsh_history`, `~/.bash_history`) and are visible to any local user via `ps aux`/`/proc/<pid>/cmdline` for the process lifetime.

**Failure scenario:** an operator runs `swctl auth create-admin ops@x.com S3cret!` on a shared jump host or CI runner; another user on the same box captures the password via `ps` or shell history, or it ends up in CI job logs.

**Remediation:** support prompting for passwords (masked stdin) as the primary path, keep flags/positional args as opt-in/discouraged. Severity: High.

Added `internal/prompt` (new package, backed by `golang.org/x/term`): when a password positional arg or the `--password`/`-p`/`AUTH_PASSWORD` flag is omitted, swctl now prompts for masked input on the real terminal instead of erroring or accepting it another way; fails fast with a clear message if stdin isn't an interactive TTY (CI/pipes), rather than hanging. Positional/flag/env-var password input still works unchanged when supplied — no breaking change for existing scripts. Wired into all 15 commands that take the authenticating `--password` flag, plus the three positional-password commands named above. `internal/prompt/password_test.go` added (6 tests, no prior test coverage existed in this repo). Verified manually: omitted password prompts; `< /dev/null` fails immediately without hanging; supplying the password positionally/via flag is unchanged. `go build`, `go vet`, and `go test ./... -timeout 10s` all pass. Committed on branch `fix/password-prompt`.

### 3. No confirmation before `auth delete-service-client`

`internal/cmd/auth/delete_service_client.go` deletes a service client immediately on a single command with no `--yes`/confirmation/dry-run, given only a `clientId` positional arg.

**Failure scenario:** a mistyped or copy-pasted wrong client ID instantly revokes a production service's ability to authenticate (e.g. `swayrider-api`'s own client), causing an outage with no undo. Severity is bounded (no user/DB data loss) but this is the closest thing to a destructive command in the tool and has zero safeguard. Severity: Medium.

### 4. gRPC transport is unencrypted (insecure credentials)

`grpcclients/internal/client/client.go` dials with `grpc.WithTransportCredentials(insecure.NewCredentials())`. Since `--auth-host`/`AUTH_HOST` is a user-supplied value, admin login credentials and issued access tokens travel in plaintext if swctl is ever pointed at a host outside the trusted internal network. This is inherited from the shared `grpcclients` module (see [[grpcclients]] review) and consistent with the platform's stated internal-only port convention, so likely an accepted design tradeoff rather than an swctl-specific bug — flagged for awareness only. Severity: Low/Info.
