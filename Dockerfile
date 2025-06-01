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
# Etapa 2: Imagen mínima (scratch)
# ------------------------------------------------------------------------------
FROM scratch

# Omitir certificados si tu servicio Go no hace HTTPS outbound.
# Para HTTPS: en builder instalar ca-certificates y copiar aquí
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER 1000:1000
WORKDIR /app

COPY --from=builder /go-sensor-data-collector /go-sensor-data-collector

EXPOSE 3000

ENTRYPOINT ["/go-sensor-data-collector"]
