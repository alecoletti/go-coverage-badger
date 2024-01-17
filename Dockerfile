FROM --platform=${BUILDPLATFORM} golang:1.21 AS builder
WORKDIR /src
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o /out/go-badger main.go
WORKDIR /bin

org.opencontainers.image.source = "https://github.com/alecoletti/go-coverage-badger"
### PRODUCTION
FROM scratch as production
COPY --from=builder /out/go-badger /bin/go-badger

ENTRYPOINT ["/bin/go-badger"]