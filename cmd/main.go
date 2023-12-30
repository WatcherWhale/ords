package main

import (
	"github.com/rs/zerolog/log"
	"github.com/watcherwhale/ords/internal/http"
)

var	version string = "dev"

func main() {
	err := http.CreateServer(version, log.Logger).Run("0.0.0.0:8080")

	if err != nil {
		panic(err)
	}
}
