package http

import (
	"net/http"

	ginlogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	v1 "github.com/watcherwhale/ords/internal/api/v1"
)

func CreateServer(version string, logger zerolog.Logger) *gin.Engine {
	r := gin.New()

	if version != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(ginlogger.SetLogger(ginlogger.WithLogger(func(ctx *gin.Context, l zerolog.Logger) zerolog.Logger {
		return logger
	})))

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"name": "oras-download",
			"version": version,
		})
	})

	v1.ConfigureRouter(r)

	return r
}
