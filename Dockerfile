FROM golang:alpine
ARG WAECHTER_VERSION
ENTRYPOINT ["/waechter"]
COPY waechter /
COPY locales /locales
COPY LICENSE /