package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

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

	width, height, imgURL, err := parseFillParams(r.URL.Path)
	if err != nil {
		ph.server.logger.Error(fmt.Sprintf("Failed to parse parameters from path: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imgData := image.ImgData{
		ImageURL: imgURL,
		Width:    width,
		Height:   height,
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

func parseFillParams(urlPath string) (width, height int, imageURL string, err error) {
	const prefix = "/fill/"
	if !strings.HasPrefix(urlPath, prefix) {
		return 0, 0, "", fmt.Errorf("invalid URL path format: %q", urlPath)
	}

	parts := strings.SplitN(urlPath[len(prefix):], "/", 4)
	if len(parts) != 4 {
		return 0, 0, "", fmt.Errorf("invalid URL path format: %q", urlPath)
	}

	width, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, "", fmt.Errorf("invalid width parameter: %w", err)
	}

	height, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, "", fmt.Errorf("invalid height parameter: %w", err)
	}

	if width <= 0 || height <= 0 {
		return 0, 0, "", fmt.Errorf("width and height must be positive integers")
	}

	imageURL = "http://" + path.Join(parts[2:]...)

	if _, err = url.ParseRequestURI(imageURL); err != nil {
		return 0, 0, "", fmt.Errorf("invalid image URL: %w", err)
	}

	return width, height, imageURL, nil
}
