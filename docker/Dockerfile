FROM golang:1.23-bullseye AS builder

RUN apt-get update && apt-get install -y gcc g++ make git && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# Production build
FROM builder AS prod-builder
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/bin/analytics ./cmd/analytics/main.go

# Development stage
FROM builder AS development
RUN apt-get update && apt-get install -y gosu && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /go/pkg/mod /app/export /app/bin
COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["air", "-c", ".air.toml"]