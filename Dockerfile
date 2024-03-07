# Etapa de construcción
FROM golang:1.22 AS builder

# Instalar swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el archivo go.mod y go.sum y descargar las dependencias
# Esto aprovecha la caché de las capas de Docker si los archivos mod/sum no se modifican
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente del proyecto
COPY . .

# Generar la documentación de Swagger
RUN swag init

# Construir la aplicación Go amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Etapa de ejecución
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario y la documentación de Swagger desde la etapa de construcción
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

# Exponer el puerto en el que tu aplicación escucha
EXPOSE 8080

CMD ["./main"]
