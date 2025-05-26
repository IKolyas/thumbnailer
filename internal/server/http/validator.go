package http

import (
	"fmt"
	"regexp"
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

	fmt.Println(urlPath)

	re := regexp.MustCompile(`^(?P<width>\d+)/(?P<height>\d+)/(?P<url>.+)$`)

	matches := re.FindStringSubmatch(urlPath[len(prefix):])
	if matches == nil {
		return fmt.Errorf("invalid URL path format: %q", urlPath)
	}

	params := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(matches) {
			params[name] = matches[i]
		}
	}

	w := params["width"]
	h := params["height"]
	rawURL := params["url"]

	width, err := strconv.Atoi(w)
	if err != nil {
		return fmt.Errorf("invalid image width: %w", err)
	}

	height, err := strconv.Atoi(h)
	if err != nil {
		return fmt.Errorf("invalid image height: %w", err)
	}

	if strings.HasPrefix(rawURL, "http:/") && !strings.HasPrefix(rawURL, "http://") {
		rawURL = strings.Replace(rawURL, "http:/", "http://", 1)
	}

	if strings.HasPrefix(rawURL, "https:/") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = strings.Replace(rawURL, "https:/", "https://", 1)
	}

	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}

	f.ImageURL = rawURL
	f.Width = width
	f.Height = height

	return nil
}
