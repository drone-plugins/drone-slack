FROM alpine:3.20
ADD drone-slack woodpecker-slack
CMD ["/woodpecker-slack"]
