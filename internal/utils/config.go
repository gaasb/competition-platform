package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var config = Config{}

type Config struct {
}

func (c *Config) setTestMode()       {}
func (c *Config) setProductionMode() {}

func init() {
	err := godotenv.Load(".env")
	HandleError(err)
}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func UpdateConfigMode() {
	if gin.Mode() == gin.ReleaseMode {
		config.setProductionMode()
	} else {
		config.setTestMode()
	}
}

func UpdateCorsSettings() {}
