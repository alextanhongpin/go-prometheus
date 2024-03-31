FROM golang:1.22.1-alpine3.19 AS builder

# Add timezone information.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser.
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

WORKDIR $GOPATH/src/alextanhongpin/app/

COPY go.mod go.mod

RUN go mod download
RUN go mod verify

COPY . .

# https://pkg.go.dev/cmd/go#hdr-Build_and_test_caching
ENV GOCACHE=/root/.cache/go-build \
		GODEBUG=gocachetest=1
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/app

FROM gcr.io/distroless/static-debian11

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /go/bin/app /go/bin/app

USER appuser:appuser

CMD ["/go/bin/app"]
