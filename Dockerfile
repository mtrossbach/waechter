FROM golang:alpine

ENTRYPOINT ["/waechter"]
COPY waechter /
COPY locales /locales
COPY LICENSE /