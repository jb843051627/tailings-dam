FROM golang:1.22-bookworm

WORKDIR /app

ENV GOTOOLCHAIN=local
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o tailings-dam .

EXPOSE 8080

ENV PORT=8080
ENV DB_PATH=/data/tailings-dam.db

VOLUME ["/data"]

CMD ["./tailings-dam"]
