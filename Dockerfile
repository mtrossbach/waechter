FROM scratch

ENTRYPOINT ["/waechter"]
COPY waechter /
COPY locales /
COPY LICENSE /