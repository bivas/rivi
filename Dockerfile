FROM alpine

RUN apk add --no-cache ca-certificates

EXPOSE 8080

VOLUME /config

COPY rivi /usr/local/bin/rivi

CMD [ "rivi", "-h" ]