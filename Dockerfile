FROM alpine:3.7
ADD ./drone-passwordstate .

ENTRYPOINT [ "./drone-passwordstate" ]
