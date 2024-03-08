# Etapa de construcción
FROM golang:1.22 as builder

# Instalar swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el archivo go.mod y go.sum y descargar las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente del proyecto
COPY . .

# Generar la documentación de Swagger
RUN swag init

# Construir la aplicación Go amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

# Etapa final
FROM alpine
COPY --from=builder /main /main
ENTRYPOINT [ "/main" ]
