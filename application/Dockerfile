FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY internal/worker/main.go .

RUN go build -o renderer main.go

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/renderer .

CMD ["./renderer"]