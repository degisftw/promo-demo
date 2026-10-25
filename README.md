# promo-demo

Deliberately small HTTP app used to exercise the Kargo promotion pipeline in
[talos-gitops](https://github.com/degisftw/talos-gitops): it reports its build
version and deployed environment, so a non-prod → prod promotion is visible
with two `curl`s.

## Endpoints

| Path | Purpose |
|---|---|
| `GET /` | `{"app":"promo-demo","version":"vX.Y.Z","env":"<APP_ENV>"}` |
| `GET /healthz` | `200 ok` — Deployment probe target |

`APP_ENV` is injected per deployment environment (`non-prod` / `prod`).
`version` is stamped at build time via `-ldflags -X main.version=`.

## Release

Push a `vX.Y.Z` tag. GitHub Actions builds `linux/arm64,linux/amd64` and pushes
`ghcr.io/degisftw/promo-demo:vX.Y.Z` (+ `:latest`). The Kargo Warehouse in
talos-gitops watches this repo's tags as Freight — new semver tags roll into
`non-prod` automatically and wait for a manual promotion to `prod`.

## Local

```sh
go test ./...
go run .                                  # :8080, version=dev
docker build -t promo-demo --build-arg VERSION=v0.0.0 .
```