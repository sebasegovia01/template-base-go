package utils

import (
	"fmt"
	"os"
)

type IEnvironment interface {
	GetEnvVar(varName string) string
	MustGetEnvVar(varName string) string
	LookupEnv(varName string) (string, bool)
	IsLambda() bool
}

type Environment struct {
	logger ILogger
}

func NewEnvironment() IEnvironment {
	return &Environment{}
}

// MustGetEnVar tries to get an env Var. It will panic if not found
func (env *Environment) MustGetEnvVar(varName string) string {
	v := os.Getenv(varName)
	if v == "" {
		message := fmt.Sprintf("%s environment variable not set.", varName)
		env.logger.Panic(nil, message, nil)
	}
	return v
}

// GetEnvVar gets and env var
func (env *Environment) GetEnvVar(varName string) string {
	return os.Getenv(varName)
}

// LookupEnv tries to get an env var. It returns a true or false if it found it
func (env *Environment) LookupEnv(varName string) (string, bool) {
	return os.LookupEnv(varName)
}

func (env *Environment) IsLambda() bool {
	_, isLambda := env.LookupEnv("AWS_LAMBDA_RUNTIME_API")

	return isLambda
}
