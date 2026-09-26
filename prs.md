# Security pull request ideas

This backlog is based on a repository-level security review of the Go API, Svelte applications, event streaming, authentication, deployment packaging, and GitHub Actions. It is ordered by security impact first, then by how safely the change can be reviewed and deployed.

The review did not include the live Cognito user pool, production host, firewall, GitHub environment protections, or a penetration test against a deployed instance. Items involving those systems need coordination with a maintainer.

## How to use this backlog

Each numbered section is a standalone implementation brief. Start a fresh task with:

> Implement PR N from `prs.md`. Read `AGENTS.md`, inspect the cited code and tests, keep the change limited to that PR's scope, run the listed tests plus the repository-required checks, and report any deployment or maintainer action that cannot be completed in the repository.

Replace `N` with the PR number and paste that PR's section if the new task cannot access this checkout. Respect the listed dependencies: some changes require Cognito, database, host, or secret rotation work that a code-only task cannot safely perform.

## Review verification

- `go test ./...`, `go vet ./...`, and `golangci-lint run` passed for the backend.
- `govulncheck ./...` found six reachable Go standard-library advisories in the Go 1.26.5 toolchain used for the review; the reported fixes are in Go 1.26.6.
- `pnpm audit --prod` reported no known production dependency advisories.
- The full npm audit reported 13 high and 7 moderate advisories in development or CI dependencies.
- A YAML parser rejected the current Dependabot configuration.
- A full frontend check was not completed because the exact package-manager invocation attempted a dependency download in the restricted review environment.

These results are a snapshot from 2026-09-01, not a guarantee that later dependency or deployment state is safe.

## Priority guide

- **P0:** Address first. Direct credential disclosure, unauthorized data access, remotely reachable denial of service, or a known vulnerable production runtime.
- **P1:** Address next. Authentication boundary weaknesses, privilege-management flaws, CI execution risks, or broken security maintenance.
- **P2:** Defense in depth. Reduces blast radius or prevents future regressions.

## Recommended order

1. Redact registration codes from request logs.
2. Upgrade and rebuild with the patched Go toolchain.
3. Prevent workflow script injection.
4. Authorize private contender event streams.
5. Bound HTTP requests and event-stream resources.
6. Validate all Cognito access-token claims.
7. Make organizer invites authorized and single-use.
8. Restore dependency security updates and scanning.
9. Replace the browser-embedded Cognito client secret.
10. Strengthen registration-code generation and guessing defenses.
11. Remove secrets from the package and deployment command line.
12. Complete the P2 hardening items.

---

## P0 — urgent

### 1. `fix: redact registration codes from request logs`

**Effort:** Small  
**Dependencies:** None

Registration codes are bearer credentials that permit scorecard changes. The lookup route places the credential in the URL, and the access logger records the full URL path.

Evidence:

- [`backend/internal/handlers/rest/contender_hdlr.go`](backend/internal/handlers/rest/contender_hdlr.go#L34) exposes `GET /codes/{registrationCode}/contender`.
- [`web/packages/lib/src/Api.ts`](web/packages/lib/src/Api.ts#L108) puts the code into that route.
- [`backend/cmd/api/main.go`](backend/cmd/api/main.go#L497) records `r.URL.Path` verbatim.

Scope:

- Immediately redact sensitive route segments in access logs.
- Prefer logging the matched route pattern, such as `/api/codes/{registrationCode}/contender`, rather than the concrete path.
- Audit other future path or query parameters before logging them.
- Document that existing application and proxy logs may already contain active credentials.

Tests:

- Send a request containing a recognizable fake registration code.
- Assert that the code does not appear in captured log output.
- Assert that method, route pattern, status, duration, and user agent are still logged.

Follow-up: remove the credential-bearing route entirely in PR 10.

### 2. `chore: require go 1.26.6 or newer`

**Effort:** Small  
**Dependencies:** None

`govulncheck` run with Go 1.26.5 reported six reachable standard-library advisories. The most relevant server-side issue is the TLS post-handshake resource-exhaustion vulnerability GO-2026-6090. All reported standard-library findings are fixed in Go 1.26.6.

Evidence:

- [`backend/go.mod`](backend/go.mod#L3) selects the Go 1.26 language version without pinning a patched toolchain.
- The review's `govulncheck ./...` result reported reachable calls for GO-2026-6218, GO-2026-6090, GO-2026-6089, GO-2026-6088, GO-2026-5972, and GO-2026-5026.

Scope:

- Set the minimum or selected toolchain to Go 1.26.6 or newer.
- Ensure GitHub Actions and release builds use that patched version.
- Rebuild and redeploy production artifacts; changing the source declaration does not patch already-built binaries.
- Update the documented prerequisite in `AGENTS.md` if the version is pinned there.

Tests:

- Record `go version` in the build job.
- Run `go test ./...` and `govulncheck ./...`.
- Confirm the built binary reports a patched Go version using `go version -m`.

Relevant advisories: [GO-2026-6090](https://pkg.go.dev/vuln/GO-2026-6090), [GO-2026-6089](https://pkg.go.dev/vuln/GO-2026-6089), [GO-2026-6088](https://pkg.go.dev/vuln/GO-2026-6088), [GO-2026-5972](https://pkg.go.dev/vuln/GO-2026-5972), [GO-2026-6218](https://pkg.go.dev/vuln/GO-2026-6218), and [GO-2026-5026](https://pkg.go.dev/vuln/GO-2026-5026).

### 3. `fix: authorize contender event subscriptions`

**Effort:** Medium  
**Dependencies:** Coordinate the frontend and backend changes in one PR

The private contender SSE endpoint accepts an enumerable contender ID without checking authentication or ownership. It streams attempts, point values, score changes, and raffle events. The public scoreboard exposes contender IDs, so targets are easy to identify.

Evidence:

- [`backend/internal/handlers/rest/events_hdlr.go`](backend/internal/handlers/rest/events_hdlr.go#L69) loads the contender and subscribes without calling the authorizer.
- [`backend/internal/handlers/rest/events_hdlr_test.go`](backend/internal/handlers/rest/events_hdlr_test.go#L64) verifies that an anonymous request receives `200 OK`.
- [`web/scorecard/src/pages/Scorecard.svelte`](web/scorecard/src/pages/Scorecard.svelte#L244) uses native `EventSource`, which cannot attach the existing `Authorization` header.

Scope:

- Require the same contender ownership check used by protected scorecard endpoints.
- Replace native `EventSource` with authenticated `fetch()` streaming, or issue a short-lived, single-purpose event ticket.
- If an event ticket is used, make it short-lived, audience-bound, single-purpose, and redacted from every log.
- Keep `/contests/{contestID}/events` public, but explicitly document its public event types.
- Return an authentication error before allocating a broker subscription.

Tests:

- Anonymous private subscription is rejected.
- A contender cannot subscribe to another contender.
- The correct contender, organizer owner, and admin can subscribe.
- Public contest events remain available anonymously.
- No rejected request allocates or leaks a broker subscription.

### 4. `fix: bound http requests and event streams`

**Effort:** Medium  
**Dependencies:** None; coordinate timeout values with SSE behavior

All HTTP timeouts are explicitly disabled, JSON bodies are decoded without a size limit, and anonymous SSE connections have no global or per-client cap. This permits slow connections, oversized request bodies, and subscription exhaustion.

Evidence:

- [`backend/cmd/api/main.go`](backend/cmd/api/main.go#L192) sets `ReadTimeout`, `ReadHeaderTimeout`, `WriteTimeout`, and `IdleTimeout` to zero.
- REST handlers call `json.NewDecoder(r.Body).Decode(...)` directly, for example in [`contender_hdlr.go`](backend/internal/handlers/rest/contender_hdlr.go#L119).
- [`backend/internal/handlers/rest/events_hdlr.go`](backend/internal/handlers/rest/events_hdlr.go#L15) permits a 1,000-event buffer for every subscription.

Scope:

- Add conservative `ReadHeaderTimeout`, `ReadTimeout`, and `IdleTimeout` values.
- Preserve long-lived SSE by leaving the global write timeout compatible with streaming and applying per-write deadlines where practical.
- Create one JSON request helper using `http.MaxBytesReader`.
- Reject unknown fields, trailing JSON values, incorrect content types, and oversized bodies consistently.
- Add global and per-client SSE connection limits and a smaller justified buffer.
- Ensure disconnects stop tickers, event readers, and goroutines promptly.

Tests:

- Oversized bodies return `413 Request Entity Too Large`.
- Slow or incomplete headers are terminated by an integration test with short test-only timeouts.
- Unknown fields and multiple JSON values return `400 Bad Request`.
- Exceeding the subscription limit is rejected without allocating a broker subscription.
- Existing SSE reconnect and keep-alive tests continue to pass.

---

## P1 — important

### 5. `ci: prevent workflow script injection`

**Effort:** Small  
**Dependencies:** None

The pull-request title and merge-group commit message are interpolated directly into a generated shell script. A crafted title can terminate the assignment and execute commands on the Actions runner. This is the same unsafe pattern documented by GitHub.

Evidence:

- [`.github/workflows/pr-title.yml`](.github/workflows/pr-title.yml#L14) embeds `${{ github.event.pull_request.title }}` and `${{ github.event.merge_group.head_commit.message }}` inside `run`.

Scope:

- Pass untrusted context values through step environment variables.
- Keep all variable expansions quoted.
- Preserve the existing lowercase and Conventional Commits checks.

Tests:

- Validate ordinary valid and invalid titles.
- Include titles containing quotes, command substitutions, newlines, glob characters, and shell metacharacters.
- Assert that malicious-looking titles are treated only as data.

Reference: [GitHub Actions script injections](https://docs.github.com/en/actions/concepts/security/script-injections).

### 6. `fix: validate cognito access token claims`

**Effort:** Medium  
**Dependencies:** Maintainer must confirm the expected issuer and app-client ID

The JWT decoder verifies RS256 signature and expiration, but it does not validate issuer, app client, or token type. This broadens the accepted credential set to unintended tokens signed by the same user pool.

Evidence:

- [`backend/internal/authorizer/jwt.go`](backend/internal/authorizer/jwt.go#L50) extracts only `username` and `exp` before accepting a token.

Scope:

- Require the exact Cognito issuer.
- Require `token_use == "access"`.
- Require `client_id` to match an explicit allowlist.
- Validate required time claims with a small documented clock skew.
- Use `sub` as the durable user identifier, with a deliberate migration plan if the database currently keys users by `username`.
- Refresh Cognito JWKS safely with caching, timeouts, key-rotation handling, and a last-known-good fallback.
- Avoid returning parsing details that help distinguish malformed tokens from other authentication failures.

Tests:

- Reject wrong issuer, wrong client, ID token, expired token, unsupported algorithm, unknown key ID, malformed claims, and bad signature.
- Accept a correctly signed access token with all expected claims.
- Test key rotation and refresh failure behavior.

Reference: [Amazon Cognito JWT verification](https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-verifying-a-jwt.html).

### 7. `fix: make organizer invites authorized and single use`

**Effort:** Medium  
**Dependencies:** Database repository change

Deleting an invite does not verify organizer ownership. Acceptance also reads the invite outside the transaction and does not check whether deletion affected a row, so concurrent requests can potentially redeem one invite more than once.

Evidence:

- [`backend/internal/usecases/organizer_ucase.go`](backend/internal/usecases/organizer_ucase.go#L173) deletes an invite without authorizing its organizer.
- [`backend/internal/usecases/organizer_ucase.go`](backend/internal/usecases/organizer_ucase.go#L187) reads before starting the acceptance transaction.

Scope:

- Resolve the invite's organizer and require organizer or admin ownership before deletion.
- Move invite lookup, expiry validation, membership creation, and consumption into one transaction.
- Lock the invite row or perform a conditional delete and require exactly one affected row.
- Make duplicate membership handling explicit and idempotent where appropriate.
- Keep the public invite-preview endpoint minimal; do not expose internal ownership or user data.

Tests:

- Anonymous and unrelated organizers cannot delete an invite.
- Owner and admin can delete it.
- An expired invite cannot be accepted.
- Two concurrent acceptance attempts result in exactly one new membership.
- A failed membership insert rolls back invite consumption.

Use `testing/synctest` for time-dependent cases, as required by the repository conventions.

### 8. `ci: restore dependency security updates`

**Effort:** Small  
**Dependencies:** None

The current Dependabot file is invalid YAML, so scheduled dependency updates are not operating. The production npm audit was clean, but the full audit reported 13 high and 7 moderate development/CI advisories.

Evidence:

- [`.github/dependabot.yml`](.github/dependabot.yml#L8) has misindented `directory` and `schedule` keys.
- A YAML parser rejects the current file.

Scope:

- Correct the YAML structure and validate it in CI.
- Cover the npm workspace, Go modules, GitHub Actions, and Docker base image.
- If the current pnpm lockfile version is unsupported by Dependabot, use Renovate or another updater rather than silently leaving it unmanaged.
- Update vulnerable development dependencies, including affected `postcss`, `undici`, `brace-expansion`, `nanoid`, `protobufjs`, and `uuid` versions.
- Add `govulncheck ./...` and `pnpm audit --prod` as required CI checks.
- Decide whether development-only audit failures block builds or are tracked with time-bounded exceptions.

Tests:

- Parse `dependabot.yml` in CI.
- Run production and full dependency audits.
- Confirm builds, checks, lint, backend tests, and E2E tests after lockfile updates.

Reference: [Configuring Dependabot version updates](https://docs.github.com/en/code-security/how-tos/secure-your-supply-chain/secure-your-dependencies/configure-version-updates).

### 9. `refactor: use cognito as a public oauth client`

**Effort:** Medium  
**Dependencies:** Requires coordinated Cognito configuration and deployment

The SPA contains a Cognito client secret and sends it from browser JavaScript. Browser applications cannot keep shared secrets and should be registered as public OAuth clients.

Evidence:

- [`web/packages/lib/src/config.json`](web/packages/lib/src/config.json#L3) contains the client secret.
- [`web/admin/src/utils/cognito.ts`](web/admin/src/utils/cognito.ts#L12) constructs a Basic authorization header with it.
- [`web/admin/src/authenticator.svelte.ts`](web/admin/src/authenticator.svelte.ts#L97) uses a 36-character UUID as the PKCE verifier; RFC 7636 specifies 43–128 characters.

Scope:

- Create a new Cognito public client without a secret.
- Generate 32 cryptographically random bytes and base64url-encode them as the PKCE verifier.
- Add a one-time `state` value and validate it on return.
- Require a verifier during every authorization-code exchange.
- Deploy the new client configuration before revoking the exposed secret.
- Remove the old secret from the current tree and, if practical, purge it from public history after rotation.
- Evaluate refresh-token rotation and a backend-for-frontend using `HttpOnly`, `Secure`, `SameSite` cookies.

Tests:

- Successful login, refresh, logout, and signup with the public client.
- Missing or incorrect verifier and state are rejected.
- No client secret appears in source files or production bundles.
- Exact redirect URIs are enforced in each environment.

References: [RFC 7636](https://datatracker.ietf.org/doc/html/rfc7636.txt) and [Cognito PKCE guidance](https://docs.aws.amazon.com/cognito/latest/developerguide/using-pkce-in-authorization-code.html).

### 10. `fix: replace registration code lookup with a session exchange`

**Effort:** Large  
**Dependencies:** PR 1 should land first

Registration codes remain in API URLs, browser history, application routes, and persistent `localStorage`. Redacting logs reduces immediate exposure but does not correct the credential lifecycle.

Evidence:

- [`web/scorecard/src/App.svelte`](web/scorecard/src/App.svelte#L39) parses codes from scorecard routes.
- [`web/scorecard/src/utils/auth.ts`](web/scorecard/src/utils/auth.ts#L44) persists complete codes.

Scope:

- Exchange a registration code once in an authenticated request body or authorization header.
- Return a scoped, revocable, expiring session credential.
- Prefer an opaque server-side session in a secure cookie; otherwise use a short-lived signed token restricted to one contender.
- Add `GET /contenders/self` so subsequent calls do not need the registration code or contender ID.
- Remove registration codes from browser routes and history.
- Define logout, expiry, revocation, and code-rotation behavior.
- Preserve the intended ability to resume up to three scorecard sessions without storing the original bearer code.

Tests:

- Registration codes never appear in request paths, logs, browser routes, or stored session objects.
- A session cannot access a different contender.
- Expired and revoked sessions fail.
- Code exchange is rate-limited and does not reveal whether failures differ internally.
- Existing resume-session and failsafe flows continue to work.

### 11. `fix: use cryptographic registration codes and throttle guessing`

**Effort:** Medium  
**Dependencies:** Coordinate with PR 10 to avoid duplicate rate-limiter work

Registration codes are credentials but are generated with `math/rand`. The eight-character space is roughly 41 bits and the unauthenticated lookup endpoint acts as an online validity oracle without rate limiting.

Evidence:

- [`backend/cmd/api/main.go`](backend/cmd/api/main.go#L52) generates codes with `math/rand.Intn`.

Scope:

- Generate codes with `crypto/rand` and unbiased alphabet selection.
- Retry database uniqueness conflicts without failing the entire ticket batch.
- Increase entropy if manual-entry usability permits.
- Add global and per-source throttles for failed exchanges, with bounded memory and expiry.
- Do not trust arbitrary `X-Real-IP`; use `RemoteAddr` unless requests are known to come through a configured trusted proxy.
- Add metrics and alerting for sustained failed-code traffic without logging submitted codes.
- If codes are no longer needed after exchange, make them one-time or rotate them after use.

Tests:

- Deterministic generator tests through an injected randomness interface.
- Distribution and alphabet tests without asserting exact random output.
- Collision retry test.
- Rate-limit, recovery-window, proxy-header, and bounded-state tests.

### 12. `fix: remove deployment secrets from command lines`

**Effort:** Medium  
**Dependencies:** Maintainer access to deployment configuration

The Debian package ships a trivial database password. Deployment replaces it by interpolating a secret into a remote shell and `sed` command, which exposes it to process arguments and mishandles shell or replacement metacharacters.

Evidence:

- [`packageroot/etc/systemd/system/climblive.service.d/override.conf`](packageroot/etc/systemd/system/climblive.service.d/override.conf#L1) contains the packaged default.
- [`.github/actions/deploy/action.yaml`](.github/actions/deploy/action.yaml#L21) inserts the secret into a remote command.

Scope:

- Remove the usable default password from the package.
- Make startup fail clearly when no database credential is configured.
- Use systemd `LoadCredential`, or a root-readable environment file installed without putting the value in process arguments.
- Add `DB_PASSWORD_FILE` support if file-based credentials are used.
- Rotate the production database password after the new mechanism is deployed.
- Quote and validate non-secret deployment inputs such as host, user, socket path, and package filename.

Tests:

- Package installation does not contain a usable database password.
- Startup succeeds with a credential file and fails without one.
- Passwords containing spaces and shell or `sed` metacharacters work unchanged.
- The secret does not appear in workflow logs, process arguments, or the final package.

---

## P2 — defense in depth

### 13. `fix: apply security headers to every response`

**Effort:** Small  
**Dependencies:** None

Security headers are applied mainly when serving an SPA index. API responses, errors, and static assets do not receive a consistent baseline, and HSTS is absent.

Evidence:

- [`backend/internal/handlers/rest/static_hdlr.go`](backend/internal/handlers/rest/static_hdlr.go#L64) applies the current header set only through the index handler.
- [`backend/cmd/api/main.go`](backend/cmd/api/main.go#L184) mounts API and application handlers without a global security-header middleware.

Scope:

- Add middleware covering API, static, error, and redirect responses.
- Send HSTS only when HTTPS is active, with a staged max-age before considering `includeSubDomains` or preload.
- Apply `X-Content-Type-Options: nosniff`, an explicit `Referrer-Policy`, and a minimal `Permissions-Policy` globally.
- Preserve the current CSP and its `frame-ancestors 'none'` protection.
- Add automated header checks for all applications and representative API responses.

Tests:

- Assert the baseline headers on successful API, API error, static asset, SPA fallback, redirect, and `404` responses.
- Assert HSTS is present in TLS mode and absent in explicitly supported plaintext development mode.
- Run each app's navigation and asset-loading tests to catch CSP regressions.

### 14. `fix: restrict production cors origins`

**Effort:** Small  
**Dependencies:** Confirm supported development origins

The API responds with `Access-Control-Allow-Origin: *` and permits the `Authorization` header from every origin. This is not currently a cookie-based authentication bypass, but it unnecessarily permits any website to call the API with credentials it obtains.

Evidence:

- [`backend/cmd/api/main.go`](backend/cmd/api/main.go#L72) grants wildcard origins in the preflight handler.
- [`backend/internal/handlers/rest/cors.go`](backend/internal/handlers/rest/cors.go#L7) also adds the wildcard origin to every API response.

Scope:

- Remove CORS in same-origin production deployments, or use an explicit environment-configured allowlist.
- Keep only the local development origins that are actually required.
- Add `Vary: Origin` when reflecting an allowed origin.
- Reject disallowed preflight requests and minimize allowed methods and headers.
- Do not enable credentialed CORS unless a concrete cookie-based design requires it.

Tests:

- Allowed production and development origins receive the expected headers.
- Unknown, malformed, and `null` origins do not receive an allow-origin header.
- Preflight tests cover allowed and disallowed method/header combinations.
- Same-origin API behavior remains unchanged when no `Origin` header is sent.

### 15. `chore: sandbox the systemd service`

**Effort:** Medium  
**Dependencies:** Test on the actual deployment host

The process drops privileges itself, but the systemd unit does not apply operating-system sandboxing or resource limits.

Evidence:

- [`packageroot/etc/systemd/system/climblive.service`](packageroot/etc/systemd/system/climblive.service#L7) defines only the executable and restart behavior under `[Service]`.

Scope:

- Prefer `User=climblive` and `Group=climblive` in systemd.
- Use `CAP_NET_BIND_SERVICE` or socket activation instead of starting the application as root solely to bind port 443.
- Add appropriate `NoNewPrivileges`, `ProtectSystem`, `ProtectHome`, `PrivateTmp`, `PrivateDevices`, `RestrictSUIDSGID`, `LockPersonality`, capability bounds, and address-family restrictions.
- Add file-descriptor, memory, task, and restart limits appropriate for expected SSE traffic.
- Ensure certificate and credential permissions remain minimal.

Tests:

- Run `systemd-analyze security climblive.service` before and after.
- Exercise TLS, database access, migrations, static serving, SSE, and graceful shutdown under the hardened unit.

### 16. `fix: support verified tls for remote mysql`

**Effort:** Medium  
**Dependencies:** Database deployment details

The MySQL DSN does not configure TLS. Production currently points to localhost, which reduces immediate exposure, but changing `DB_HOST` to a remote database would silently send credentials and data without transport verification.

Evidence:

- [`backend/cmd/api/main.go`](backend/cmd/api/main.go#L117) passes host and credentials to the repository without a TLS policy.
- [`backend/internal/repository/db.go`](backend/internal/repository/db.go#L20) constructs a plaintext-capable DSN without registering or requiring a verified TLS configuration.

Scope:

- Require verified TLS for non-loopback database hosts.
- Support a configured CA bundle and server name.
- Fail closed for remote hosts when TLS is absent or invalid.
- Document an explicit local-only exception rather than silently selecting plaintext.

Tests:

- Verified TLS succeeds.
- Wrong CA, wrong hostname, and plaintext remote configuration fail.
- The documented loopback development configuration continues to work.

### 17. `ci: pin third party actions by commit sha`

**Effort:** Small  
**Dependencies:** PR 8 should enable automated update PRs

Workflow dependencies currently use mutable major-version tags. A compromised upstream tag could affect builds or deployment.

Evidence:

- [`.github/workflows/pull-request.yml`](.github/workflows/pull-request.yml#L111) uses mutable tags for checkout, setup, artifact, lint, and SQLC actions.
- [`.github/workflows/release.yml`](.github/workflows/release.yml#L77) uses a mutable release-action tag in the privileged deployment workflow.

Scope:

- Pin every third-party action to a reviewed full commit SHA.
- Retain the human-readable release version in a comment.
- Enable Dependabot or Renovate updates for GitHub Actions.
- Give every job only the permissions it requires.
- Review whether pull-request deployment to the shared test environment should require environment approval.

Tests:

- Parse every workflow and assert that external `uses:` references are full 40-character SHAs.
- Run the pull-request and release workflows from a safe test branch or dry-run environment.
- Confirm the updater recognizes and proposes new action SHAs without widening permissions.

### 18. `fix: align api validation with database limits`

**Effort:** Medium  
**Dependencies:** None

Several validators check presence or shape but not the database's maximum string lengths. Oversized values therefore consume unnecessary memory and fail as repository errors rather than controlled validation errors.

Evidence:

- [`backend/internal/usecases/validators/contest.go`](backend/internal/usecases/validators/contest.go#L34) validates required fields but not persisted string lengths.
- [`backend/internal/usecases/validators/tick.go`](backend/internal/usecases/validators/tick.go#L14) compares attempt ordering but permits negative and arbitrarily large counts.
- [`backend/database/climblive.sql`](backend/database/climblive.sql#L19) defines the actual `VARCHAR` limits that API validation should mirror.

Scope:

- Add rune or byte limits matching every persisted name, description, location, country, color, and info column.
- Add sensible maximum attempt counts and reject negative attempt values.
- Reject non-positive resource IDs at the HTTP boundary.
- Keep HTML sanitization for contest info and add bypass-regression cases.
- Return consistent `400 Bad Request` responses for validation failures.

Tests:

- Boundary tests at maximum minus one, maximum, and maximum plus one.
- Unicode tests that distinguish bytes from runes where relevant.
- Negative and extreme numeric input tests.
- Sanitizer tests for scripts, event handlers, unsafe URLs, malformed markup, and encoded variants.

### 19. `test: add an api authorization matrix`

**Effort:** Medium  
**Dependencies:** Best after PRs 3, 6, and 7

Authorization is spread across handlers and use cases. Existing unit tests are substantial, but there is no single regression suite proving the policy for every route and role.

Evidence:

- [`backend/cmd/api/main.go`](backend/cmd/api/main.go#L271) assembles all use cases and handlers without a central route-policy declaration.
- [`backend/internal/handlers/rest/contender_hdlr.go`](backend/internal/handlers/rest/contender_hdlr.go#L32) illustrates route registration being separate from authorization inside use cases.

Scope:

- Build a table covering every method and route as anonymous, contender owner, different contender, organizer owner, different organizer, and admin.
- Explicitly classify public endpoints.
- Verify both response status and absence of side effects.
- Add tests ensuring authentication is checked before expensive allocation or mutation.
- Include invite, SSE, archive/restore, transfer, raffle, score-engine, scrub, and results-download routes.

The matrix should become the authoritative documentation for the API's access-control policy.

Tests:

- Make the table itself executable as integration or handler tests.
- Require every registered route to have a matching policy row so newly added routes fail the suite until classified.
- Include representative not-found resources to ensure authorization checks do not leak resource existence.

### 20. `docs: add a security policy and threat model`

**Effort:** Small  
**Dependencies:** None

A public scoring platform benefits from clear disclosure and trust assumptions, especially because contenders self-report results and registration codes are physical bearer credentials.

Evidence:

- The repository currently has no `SECURITY.md` or equivalent vulnerability-disclosure policy.
- The repository also has no root project overview that defines supported releases, a reporting channel, or credential trust boundaries.

Scope:

- Add `SECURITY.md` with supported versions, a private reporting channel, expected response times, and safe-harbor language.
- Document assets, actors, trust boundaries, and intentionally public data.
- State the security properties and lifecycle of organizer tokens, contender credentials, and invites.
- Document operational responsibilities for Cognito, database access, TLS keys, logs, backups, and dependency alerts.
- Include a release checklist covering vulnerability scans, secret scanning, and production rollback.

Tests:

- Have the maintainer verify the reporting address and supported-version policy before merging.
- Check that all internal links and contact details render correctly on the public repository.
- Add a lightweight CI link check if the repository already uses documentation validation.

---

## Review findings that do not currently need a security PR

- Repository SQL uses SQLC-generated parameterized queries; no SQL injection path was identified.
- Contest rich text is sanitized with Bluemonday before the scorecard renders it as HTML.
- The application CSP prevents third-party script execution and framing.
- API responses use `Cache-Control: no-store`.
- TLS is configured with a minimum version of TLS 1.2.
- The server drops privileges before accepting requests.
- Production npm dependencies had no known advisories at the time of review.

These controls should be preserved with regression tests when adjacent code changes.
