# Template Base Go

Este proyecto es una plantilla base para aplicaciones en Go, diseñada para facilitar el desarrollo, pruebas y despliegue de aplicaciones robustas y escalables.

## Requisitos

- Go versión 1.22 o superior

## Configuración Inicial

Antes de comenzar, asegúrate de tener configurado Go correctamente y de tener todas las herramientas necesarias instaladas, incluyendo `swag` para la generación de documentación de Swagger y `air` para la recarga automática del servidor durante el desarrollo.

### Instalar Dependencias
Para descargar e instalar las dependencias necesarias:

`go mod download`

### Cambiar Nombres de Módulos

Si necesitas cambiar el nombre del módulo de tu proyecto, puedes utilizar el script `change-base-name.sh`:

`sh ./change-base-name.sh <nombre-actual> <nuevo-nombre>`

### Ejecutar el Servidor
Para iniciar el servidor de desarrollo:

`go run main.go`

### Inicializar Configuración para Recarga Automática
Para preparar tu proyecto para utilizar air y beneficiarte de la recarga automática durante el desarrollo, primero ejecuta:

`air init`

Luego, para iniciar el servidor con recarga automática:

`air`

### Generar Documentación con Swagger

`swag init`

### Compilar el Código
Para compilar tu proyecto:

`go build`

### Ejecutar Pruebas
Para ejecutar las pruebas de tu proyecto:

`go test ./...`
