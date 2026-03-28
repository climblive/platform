# Plan: Deploy Service as a Container Wrapped in a Deb Package

Addresses [climblive/platform#756](https://github.com/climblive/platform/issues/756).

## Background

The current deployment packages a bare Go binary and static web assets into a
`.deb` file. The binary is executed directly by systemd. This works well but
means the runtime environment is whatever happens to be on the host.

Containerising the backend would give us a reproducible, testable deliverable
while keeping the cheap and simple deb-based deployment model.

## Approach

Only the **backend API service** is containerised. The frontend apps are static
files served by the host's nginx and gain nothing from running inside a
container.

The backend container image is exported as a tarball, bundled into the `.deb`
package, and loaded into Docker on the target host at install time. Systemd
manages the container lifecycle the same way it manages the binary today.

## Changes Required

### 1. Improve the backend Dockerfile

Convert to a multi-stage build so the final image is small and contains only
the static binary.

```dockerfile
FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -o ./climblive ./cmd/api/main.go

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/climblive /climblive
ENTRYPOINT ["/climblive"]
```

### 2. Update the `build-backend` action

Replace the Go toolchain setup and bare `go build` with a Docker build that
exports the image as a tarball artifact.

```yaml
steps:
  - run: |
      docker build -t climblive-api:latest backend/
      docker save climblive-api:latest -o climblive-api.tar
    shell: bash
  - uses: actions/upload-artifact@v4
    with:
      name: backend
      path: climblive-api.tar
```

### 3. Update the `package` action

Replace the binary copy (`/usr/bin/climblive`) with the image tarball placed at
`/usr/share/climblive/climblive-api.tar`.

### 4. Add `docker.io` as a deb dependency

```
Depends: nginx, certbot, python3-certbot-nginx, docker.io
```

### 5. Rewrite `postinst`

- Create `/etc/climblive/environment` with default database connection values
  (only on first install, to preserve existing config on upgrades).
- Load the container image: `docker load -i /usr/share/climblive/climblive-api.tar`.

### 6. Add `prerm` script

Stop and remove the running container before the package is removed or
upgraded.

### 7. Replace the systemd service

The service should depend on `docker.service` and run the container with host
networking so the API keeps listening on `127.0.0.1:8090` as before.

```ini
[Unit]
Description=ClimbLive 2.0
After=syslog.target docker.service
Requires=docker.service
StartLimitIntervalSec=60
StartLimitBurst=5

[Service]
ExecStartPre=-/usr/bin/docker stop climblive-api
ExecStartPre=-/usr/bin/docker rm climblive-api
ExecStart=/usr/bin/docker run --name climblive-api --rm --network host \
  --env-file /etc/climblive/environment climblive-api:latest
ExecStop=/usr/bin/docker stop climblive-api
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### 8. Remove `climblive.service.d/override.conf`

Environment variables move from a systemd override into the Docker env file at
`/etc/climblive/environment`.

### 9. Update the `deploy` action

Point the `sed` that injects the database password at
`/etc/climblive/environment` instead of the old systemd override file.

## Files Affected

| File | Change |
|------|--------|
| `backend/Dockerfile` | Multi-stage build |
| `.github/actions/build-backend/action.yaml` | Docker build + save |
| `.github/actions/package/action.yaml` | Bundle tarball instead of binary |
| `.github/actions/deploy/action.yaml` | Write password to env file |
| `packageroot/DEBIAN/control` | Add `docker.io` dependency |
| `packageroot/DEBIAN/postinst` | Load image, create env file |
| `packageroot/DEBIAN/prerm` | New — stop container |
| `packageroot/etc/systemd/system/climblive.service` | Run container |
| `packageroot/etc/systemd/system/climblive.service.d/override.conf` | Remove |

## What Stays the Same

- Frontend apps are still built, packaged, and served by the host nginx
  exactly as before.
- Nginx configuration is unchanged.
- The API still binds to `0.0.0.0:8090`; nginx still proxies `/api` there.
- The deploy workflow structure (build → package → deploy) is unchanged.
- E2E tests already build Docker images and are unaffected.

## Open Questions

- **Image size in the `.deb`**: The exported tarball will be larger than a bare
  binary. Acceptable given the benefits, but worth measuring.
- **Docker vs Podman**: Docker is chosen here for simplicity. Podman would
  offer rootless containers and better systemd integration, but requires more
  setup.
- **Image tagging**: This plan tags as `latest`. A version-tagged image
  (e.g. `climblive-api:2.4.1`) is possible but requires passing the version to
  the build step.
