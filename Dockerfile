FROM golang:alpine

ENTRYPOINT ["/waechter"]
COPY waechter /
COPY locales /
COPY LICENSE /