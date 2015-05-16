# Docker image for Drone's slack notification plugin
#
#     go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-slack .

FROM gliderlabs/alpine:3.1
RUN apk-install ca-certificates
ADD drone-slack /bin/
ENTRYPOINT ["/bin/drone-slack"]