FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./ 
RUN go mod download

COPY sfu/ ./sfu/

RUN CGO_ENABLED=0 GOOS=linux go build -o /sfu ./sfu

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /sfu .
EXPOSE 8081
CMD ["./sfu"]