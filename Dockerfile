FROM alpine:3.7
MAINTAINER Tomas Dabasinskas <tomas@dabasinskas.net>
COPY ./release/linux/amd64 .

ENTRYPOINT [ "./drone-passwordstate" ]
