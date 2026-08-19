FROM golang:1.22-alpine AS builder

ENV GOTOOLCHAIN=local
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /tailings-dam .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /tailings-dam .

EXPOSE 8080

CMD ["/app/tailings-dam"]
