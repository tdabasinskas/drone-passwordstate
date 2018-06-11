FROM alpine
ADD ./release/drone-passwordstate /bin/
RUN apk -Uuv add ca-certificates

ENTRYPOINT /bin/drone-passwordstate
