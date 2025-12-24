# ------------------------------------------------------------------------------
# Etapa 1: Build de la aplicación Go (solo SHT31)
# ------------------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

# 1) Cachear dependencias
COPY go.mod go.sum ./
RUN go mod tidy -v && go mod download

# 2) Copiar todo el código (incluye cmd/server y demás paquetes)
COPY . .

# 3) Compilar el binario estático de Go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /go-sensor-data-collector ./cmd/server

# ------------------------------------------------------------------------------
# Etapa 2: Imagen mínima (alpine para health checks)
# ------------------------------------------------------------------------------
FROM alpine:latest

# Instalar wget para health checks
RUN apk add --no-cache wget ca-certificates

# Crear usuario no-root
RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser

USER appuser
WORKDIR /app

COPY --from=builder /go-sensor-data-collector /go-sensor-data-collector

EXPOSE 3000

ENTRYPOINT ["/go-sensor-data-collector"]
