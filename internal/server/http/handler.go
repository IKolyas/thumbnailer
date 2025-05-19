package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/IKolyas/image-previewer/internal/core/image"
	"github.com/IKolyas/image-previewer/internal/storage/source"
	"github.com/davidbyttow/govips/v2/vips"
)

type PreviewerHandler struct {
	server Server
}

func (ph *PreviewerHandler) fill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//nolint
	ctx = context.WithValue(ctx, "Headers", r.Header)

	imageRequest := FillImageRequest{}

	err := imageRequest.validate(r.URL.Path)
	if err != nil {
		ph.server.logger.Error(fmt.Sprintf("Failed to parse parameters from path: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imgData := image.ImgData{
		ImageURL: imageRequest.ImageURL,
		Width:    imageRequest.Width,
		Height:   imageRequest.Height,
		Format:   vips.ImageTypeUnknown,
		Action:   image.ImageActionFill,
	}

	imageData, err := ph.server.storage.Get(ctx, &imgData)
	if err != nil {
		ph.server.logger.Error(fmt.Sprintf("Failed to get image source: %v", err))
		var sourceErr *source.Error
		if errors.As(err, &sourceErr) {
			http.Error(w, err.Error(), sourceErr.Code())
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Length", fmt.Sprint(len(imageData)))
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(imageData); err != nil {
		ph.server.logger.Error(fmt.Sprintf("Failed to write response: %v", err))
	}
}
