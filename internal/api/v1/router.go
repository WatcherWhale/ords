package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/watcherwhale/ords/internal/config"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func ConfigureRouter(r *gin.Engine) {
    router := r.Group("/v1")

    // Translate the registry name to a valid registry url
    router.Use(getRegistry)

    // Attach auth credentials if supplied
    router.Use(attachOCIAuth)

    // Register Routes
    router.GET("/download/:registry", downloadFile)
    router.GET("/info/:registry", infoFile)
}

func getRegistry(ctx *gin.Context) {
	registryParam := ctx.Param("registry")
	registry, ok := config.Config.Registries[registryParam]

	if !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "Registry does not exist",
		})

		return
	}

    ctx.Set("registry", registry)
	ctx.Next()
}

func attachOCIAuth(ctx *gin.Context) {
    username, password, ok := ctx.Request.BasicAuth()

    if ok {
        registry, _ := ctx.Get("registry")

        client := auth.Client{
            Client: retry.DefaultClient,
            Cache: auth.DefaultCache,
            Credential: auth.StaticCredential(registry.(string), auth.Credential{
                Username: username,
                Password: password,
            }),
        }

        ctx.Set("ociauth", client)
    }

	ctx.Next()
}
