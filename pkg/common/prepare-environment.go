package common

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func splitEnv(env string) [2]string {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			return [2]string{env[:i], env[i+1:]}
		}
	}
	return [2]string{env, ""}
}

func PrepareEnvironment(requiredEnvVariableNames []string) (map[string]string, error) {

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("Error loading .env file: %v", err)
		}
	}

	for _, variableName := range requiredEnvVariableNames {
		if os.Getenv(variableName) == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", variableName)
		}
	}

	envVars := make(map[string]string)
	for _, env := range os.Environ() {
		pair := splitEnv(env)
		envVars[pair[0]] = pair[1]
	}

	return envVars, nil
}
