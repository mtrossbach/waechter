FROM golang:alpine AS builder

RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates

WORKDIR /go/delivery
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /waechter github.com/mtrossbach/waechter/cmd/waechter


FROM scratch
WORKDIR /
ENV TZ=Europe/Berlin

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /waechter /waechter
COPY ./locales /locales

ENTRYPOINT ["/waechter"]