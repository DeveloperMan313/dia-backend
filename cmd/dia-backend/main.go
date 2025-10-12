package main

import (
	"dia-backend/internal/api"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// @title           Let There Be Light API
// @version         1.0
// @description     REST API для веб-приложения расчета освещения комнат

// @contact.name   API Support
// @contact.url    https://github.com/yourusername/dia-backend
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8001
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	err := godotenv.Load("deploy/.env")
	if err != nil {
		panic(err)
	}

	logrus.SetLevel(logrus.ErrorLevel)
	api.StartServer()
}
