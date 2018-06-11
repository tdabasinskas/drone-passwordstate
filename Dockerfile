FROM alpine:3.7
COPY ./release .

ENTRYPOINT [ "./drone-passwordstate" ]
