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

func downloadFile(ctx *gin.Context) {
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

    _, bytes, err := oras.FetchBytes(c, repo, ctx.Query("tag"), oras.DefaultFetchBytesOptions)

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

	logger.Info().Msgf("Downloading digest '%s'", manifest.Layers[0].Digest.String())

    _, reader, err := repo.Blobs().FetchReference(c, manifest.Layers[0].Digest.String())

    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
            "message": "blob",
        })

		logger.Warn().Err(err).Msg("Error when fetching file")

        return
    }

    defer reader.Close()

    ctx.DataFromReader(http.StatusOK, manifest.Layers[0].Size, manifest.ArtifactType, reader, map[string]string{})
}
