FROM alpine:3.6

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ADD release/linux/amd64/drone-slack /bin/
ENTRYPOINT ["/bin/drone-slack"]
