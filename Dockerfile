FROM golang as builder
WORKDIR /app
COPY . /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o github-pr-comment /app/cmd/github-pr-comments/main.go

FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM busybox
COPY docker/in docker/out docker/check /opt/resource/
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /app/github-pr-comment /opt/resource/github-pr-comment
RUN chmod +x /opt/resource/*
