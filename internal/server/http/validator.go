package http

import (
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type FillImageRequest struct {
	ImageURL string
	Width    int
	Height   int
}

func (f *FillImageRequest) validate(urlPath string) error {
	const prefix = "/fill/"
	if !strings.HasPrefix(urlPath, prefix) {
		return fmt.Errorf("invalid URL path format: %q", urlPath)
	}

	parts := strings.SplitN(urlPath[len(prefix):], "/", 4)
	if len(parts) != 4 {
		return fmt.Errorf("invalid URL path format: %q", urlPath)
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid width parameter: %w", err)
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid height parameter: %w", err)
	}

	if width <= 0 || height <= 0 {
		return fmt.Errorf("width and height must be positive integers")
	}

	imageURL := "http://" + path.Join(parts[2:]...)

	if _, err = url.ParseRequestURI(imageURL); err != nil {
		return fmt.Errorf("invalid image URL: %w", err)
	}

	f.ImageURL = imageURL
	f.Width = width
	f.Height = height

	return nil
}
