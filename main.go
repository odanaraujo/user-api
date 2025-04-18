package main

import (
	"github.com/odanaraujo/user-api/infrastructure/loggers"
	"github.com/odanaraujo/user-api/router"
)

func main() {
	defer loggers.Close()

	r := router.NewRouter()
	r.Run(":8080")
}
