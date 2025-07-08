package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/promhippie/jenkins_exporter/pkg/command"
)

func main() {
	if env := os.Getenv("JENKINS_EXPORTER_ENV_FILE"); env != "" {
		_ = godotenv.Load(env)
	}

	if err := command.Run(); err != nil {
		os.Exit(1)
	}
}
