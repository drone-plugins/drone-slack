# Docker image for Drone's slack notification plugin
#
#     docker build --rm=true -t plugins/drone-slack .

FROM library/golang:1.4

# copy the local package files to the container's workspace.
ADD . /go/src/github.com/drone-plugins/drone-slack/

# build the slack plugin inside the container.
RUN go get github.com/drone-plugins/drone-slack/... && \
    go install github.com/drone-plugins/drone-slack

# run the slack plugin when the container starts
ENTRYPOINT ["/go/bin/drone-slack"]
