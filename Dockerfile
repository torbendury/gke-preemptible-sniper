ARG GOLANG_VERSION=1.22.6-alpine

# Dev Stage
FROM golang:${GOLANG_VERSION} AS dev
WORKDIR /app
CMD ["sh"]

# Build Stage
FROM golang:${GOLANG_VERSION} AS build

ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

RUN apk update && apk add --no-cache git
WORKDIR /src
COPY . .
RUN go mod download && go mod verify
RUN mkdir /app && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/app cmd/main.go

# Release Stage
FROM scratch AS release

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

USER appuser:appuser

COPY --from=build /app/app /go/bin/gke-preemptible-sniper
ENTRYPOINT ["/go/bin/gke-preemptible-sniper"]
