#!/bin/bash

# Verifica si se proporcionaron los dos parámetros necesarios
if [ -z "$1" ] || [ -z "$2" ]; then
  echo "Error: No se definieron los parámetros requeridos."
  echo "Uso: $0 <nombreModuloAntiguo> <nombreModuloNuevo>"
  echo "Ejemplo: $0 template-base-go template-base"
  exit 1
fi

# Los nombres de módulo se toman de los argumentos pasados al script
nombreModuloAntiguo="$1"
nombreModuloNuevo="$2"

echo "Actualizando el nombre del módulo de '$nombreModuloAntiguo' a '$nombreModuloNuevo'"

# Actualiza el nombre del módulo en go.mod
sed -i '' -e "s/$nombreModuloAntiguo/$nombreModuloNuevo/g" go.mod

# Actualiza el nombre del módulo en todos los archivos .go
find . -type f -name '*.go' -exec sed -i '' -e "s/$nombreModuloAntiguo/$nombreModuloNuevo/g" {} +

echo "Actualización completada."
