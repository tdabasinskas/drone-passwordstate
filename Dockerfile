FROM alpine:3.7
COPY ./release/linux/amd64 .

ENTRYPOINT [ "./drone-passwordstate" ]
