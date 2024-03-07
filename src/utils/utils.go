package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Container struct {
	Environment IEnvironment
}

func NewContainer(env IEnvironment) *Container {
	return &Container{
		Environment: env,
	}
}

// ParseEnvironmentFile reads .env file and cretes env variables
func ParseEnvironmentFile(localEnv string) error {
	// Definir la ruta del archivo .env relativa al directorio actual de ejecución
	file := ".env"

	if localEnv != "" {
		file = fmt.Sprintf(".env%s", localEnv)
	}

	// Abre el archivo .env usando una ruta relativa
	handler, err := os.Open(file)
	if err != nil {
		log.Fatalf("Failed to open file path, %v", err)
	}
	defer handler.Close()

	// Leer y establecer las variables de entorno
	return ReadFileAndSetEnv(handler)
}

// ReadFileAndSetEnv takes a reader and will set its keys as env variables
func ReadFileAndSetEnv(handle io.Reader) error {
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") { // Ignora líneas vacías o comentadas
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("ignoring malformed line in env file: %s", line)
			continue
		}
		key := parts[0]
		value := parts[1]
		os.Setenv(key, value)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error reading env file: %v", err)
		return err
	}
	return nil
}
