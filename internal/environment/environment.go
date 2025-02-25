package environment

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not present, using process envs and defaults")
	} else {
		log.Println(".env loaded")
	}
}
