# Step 1 - Build binary
FROM golang:1.15-alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR /src/
ADD . /src/
RUN CGO_ENABLED=0 go build -o /bin/jiraalert
COPY ./.env /bin/.env


FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/jiraalert /bin/jiraalert
COPY --from=builder /bin/.env /bin/.env

WORKDIR /bin/

EXPOSE 2112

ENTRYPOINT ["/bin/jiraalert"]