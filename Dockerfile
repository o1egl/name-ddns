FROM golang:1.18-alpine as builder

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN go build -o name-ddns .

FROM alpine:3.15

COPY --from=builder /app/name-ddns /bin/name-ddns

CMD ["/bin/name-ddns"]