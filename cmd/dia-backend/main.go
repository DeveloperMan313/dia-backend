package main

import (
	"dia-backend/internal/api"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.ErrorLevel)
	api.StartServer()
}
