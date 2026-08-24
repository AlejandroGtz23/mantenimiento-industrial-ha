# Compila el servidor Go una sola vez; la imagen final no contiene el código fuente.
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY backend ./backend
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/maintenance-api ./backend/cmd

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata wget \
    && addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=build /out/maintenance-api ./maintenance-api
RUN mkdir -p /app/uploads/mantenimiento && chown -R app:app /app
USER app
EXPOSE 8081
CMD ["./maintenance-api"]
