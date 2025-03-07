FROM golang:1.18 AS builder

COPY . /app
WORKDIR /app

# RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/appbin main.go


FROM debian:stable-slim

# Following commands are for installing CA certs (for proper functioning of HTTPS and other TLS)
RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

# Add new user 'appuser'. App should be run without root privileges as a security measure
RUN adduser --home "/appuser" --disabled-password appuser \
    --gecos "appuser,-,-,-"
USER appuser

COPY --from=builder /app/appbin /home/appuser/app
WORKDIR /home/appuser/app

# Since running as a non-root user, port bindings < 1024 are not possible
# 8000 for HTTP; 8443 for HTTPS;
EXPOSE 8000

CMD ["./appbin"]