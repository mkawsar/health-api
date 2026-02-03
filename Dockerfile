# Single-stage Dockerfile for development and deployment.
# Build: docker build -t health-api .  or  docker build --build-arg ENV=development -t health-api .
# Runtime: set MODE=debug (dev) or MODE=release (deploy) via .env or compose.

ARG ENV=development

FROM golang:1.23-alpine

ARG ENV
WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Production: stripped binary. Development: faster build, no strip.
RUN set -e; \
  if [ "$ENV" = "production" ]; then \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/health-api ./main.go; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/migrate ./cmd/migrate/main.go; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/seed ./cmd/seed/main.go; \
  else \
    CGO_ENABLED=0 GOOS=linux go build -o /app/health-api ./main.go; \
    CGO_ENABLED=0 GOOS=linux go build -o /app/migrate ./cmd/migrate/main.go; \
    CGO_ENABLED=0 GOOS=linux go build -o /app/seed ./cmd/seed/main.go; \
  fi

# Default .env (override with env_file or environment in compose)
COPY .env.example .env

EXPOSE 8080

CMD ["./health-api"]
