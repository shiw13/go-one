package env

import (
	"fmt"
	"os"
)

const (
	DeployEnvDev  = "dev"
	DeployEnvTest = "test"
	DeployEnvProd = "prod"
)

var (
	DeployEnv string
)

func MustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Errorf("env: get env failed %s", key))
	}
	return val
}

func TryGetEnv(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}
	return val
}

func init() {
	DeployEnv = TryGetEnv("DEPLOY_ENV", DeployEnvProd)
}
