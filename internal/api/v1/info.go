package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func infoFile(ctx *gin.Context) {
	logger := log.Logger.With().Str("path", "/v1/download").Logger()

    registry, _ := ctx.Get("registry")
    ociauth, authenticated := ctx.Get("ociauth")


    c := context.Background()

    repo, err := remote.NewRepository(registry.(string) + "/" + ctx.Query("image"))

    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
            "message": "Could not find repository",
        })

		logger.Warn().Err(err).Msg("Error when creating repository")

        return
    }

    if authenticated {
        ociauthclient := ociauth.(auth.Client)
        repo.Client = &ociauthclient
    }

    desc, bytes, err := oras.FetchBytes(c, repo, ctx.Query("tag"), oras.DefaultFetchBytesOptions)

    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
            "message": "descriptor",
        })

		logger.Warn().Err(err).Msg("Error when fetching manifest")
        return
    }
	

	var manifest ocispec.Manifest
	if err := json.Unmarshal(bytes, &manifest); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
            "message": "manifest",
        })

		logger.Warn().Err(err).Msg("Error when parsing manifest")
	}

	ctx.JSON(http.StatusOK, gin.H{
		"image": gin.H{
			"name": registry.(string) + "/" + ctx.Query("image"),
			"tag": ctx.Query("tag"),
			"digest": desc.Digest.String(),
		},
		"file": gin.H{
			"mediaType": manifest.Config.MediaType,
			"name": manifest.Layers[0].Annotations["org.opencontainers.image.title"],
			"size": manifest.Layers[0].Size,
			"checksum": manifest.Layers[0].Digest,
		},
	})
}
