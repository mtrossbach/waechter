FROM golang:alpine
ENV WAECHTER_VERSION=$WAECHTER_VERSION
ENTRYPOINT ["/waechter"]
COPY waechter /
COPY locales /locales
COPY LICENSE /
