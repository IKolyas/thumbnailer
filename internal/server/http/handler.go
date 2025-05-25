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

const (
	headerContentLength = "Content-Length"
	headerContextKey    = "Headers"
)

type PreviewerHandler struct {
	server Server
}

func NewPreviewerHandler(server Server) *PreviewerHandler {
	return &PreviewerHandler{
		server: server,
	}
}

func (ph *PreviewerHandler) Fill(w http.ResponseWriter, r *http.Request) {
	ctx := ph.prepareContext(r)
	
	imageRequest, err := ph.parseAndValidateRequest(r)
	if err != nil {
		ph.handleError(w, "Failed to parse parameters from path", err, http.StatusBadRequest)
		return
	}

	imgData := ph.createImageData(imageRequest)
	
	imageData, err := ph.server.storage.Get(ctx, imgData)
	if err != nil {
		ph.handleStorageError(w, err)
		return
	}

	ph.writeResponse(w, imageData)
}

func (ph *PreviewerHandler) prepareContext(r *http.Request) context.Context {
	ctx := r.Context()
	return context.WithValue(ctx, headerContextKey, r.Header)
}

func (ph *PreviewerHandler) parseAndValidateRequest(r *http.Request) (*FillImageRequest, error) {
	imageRequest := &FillImageRequest{}
	if err := imageRequest.validate(r.URL.Path); err != nil {
		return nil, err
	}
	return imageRequest, nil
}

func (ph *PreviewerHandler) createImageData(req *FillImageRequest) *image.ImgData {
	return &image.ImgData{
		ImageURL: req.ImageURL,
		Width:    req.Width,
		Height:   req.Height,
		Format:   vips.ImageTypeUnknown,
		Action:   image.ImageActionFill,
	}
}

func (ph *PreviewerHandler) handleError(w http.ResponseWriter, message string, err error, statusCode int) {
	ph.server.logger.Error(fmt.Sprintf("%s: %v", message, err))
	http.Error(w, err.Error(), statusCode)
}

func (ph *PreviewerHandler) handleStorageError(w http.ResponseWriter, err error) {
	message := "Failed to get image source"
	var sourceErr *source.Error
	if errors.As(err, &sourceErr) {
		ph.handleError(w, message, err, sourceErr.Code())
		return
	}
	ph.handleError(w, message, err, http.StatusInternalServerError)
}

func (ph *PreviewerHandler) writeResponse(w http.ResponseWriter, imageData []byte) {
	w.Header().Set(headerContentLength, fmt.Sprint(len(imageData)))
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(imageData); err != nil {
		ph.server.logger.Error(fmt.Sprintf("Failed to write response: %v", err))
	}
}
