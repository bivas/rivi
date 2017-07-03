FROM debian:stretch-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 8080

VOLUME /config

COPY rivi /usr/local/bin/rivi

CMD [ "rivi", "-h" ]
