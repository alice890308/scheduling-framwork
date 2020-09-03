FROM debian:stretch-slim

WORKDIR /

COPY _output/bin/alice-scheduler /usr/local/bin

CMD ["alice-scheduler"]

