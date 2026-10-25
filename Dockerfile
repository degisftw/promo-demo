# syntax=docker/dockerfile:1

# Cross-compiles on the build host's platform (no QEMU for the Go build):
# TARGETARCH comes from buildx when building multi-arch.
FROM --platform=$BUILDPLATFORM golang:1.27.2 AS build
ARG TARGETOS=linux
ARG TARGETARCH
ARG VERSION=dev
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w -X main.version=$VERSION" -o /out/promo-demo .

# Static binary, nonroot user, no shell — passes the cluster's baseline
# Pod Security standard without any SecurityContext.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/promo-demo /promo-demo
EXPOSE 8080
ENTRYPOINT ["/promo-demo"]