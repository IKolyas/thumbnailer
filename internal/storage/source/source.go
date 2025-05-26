package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/IKolyas/image-previewer/internal/core/image"
)

type Storage interface {
	Get(ctx context.Context, imgData *image.ImgData) ([]byte, error)
}

type Error struct {
	Message    string
	StatusCode int
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Code() int {
	return e.StatusCode
}

func Get(ctx context.Context, imgData *image.ImgData) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", imgData.ImageURL, nil)
	if err != nil {
		return nil, &Error{
			Message:    fmt.Sprintf("failed to create request: %s", err),
			StatusCode: http.StatusInternalServerError,
		}
	}

	headers, ok := ctx.Value("Headers").(http.Header)
	if !ok {
		return nil, &Error{
			Message:    "invalid headers in context",
			StatusCode: http.StatusInternalServerError,
		}
	}
	req.Header = headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &Error{
			Message:    fmt.Sprintf("failed to download image: %s", err),
			StatusCode: http.StatusInternalServerError,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &Error{
			Message:    fmt.Sprintf("unexpected status code: %v", resp.StatusCode),
			StatusCode: resp.StatusCode,
		}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, &Error{
			Message:    "file is not an image",
			StatusCode: http.StatusUnsupportedMediaType,
		}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &Error{
			Message:    fmt.Sprintf("failed to read image data: %s", err),
			StatusCode: http.StatusInternalServerError,
		}
	}

	vipsImg, err := image.NewImage(data)
	if err != nil {
		return nil, &Error{
			Message:    fmt.Sprintf("failed to create vips image: %s", err),
			StatusCode: http.StatusInternalServerError,
		}
	}

	switch imgData.Action {
	case image.ImageActionFill:
		res, err := vipsImg.Fill(imgData)
		if err != nil {
			return nil, &Error{
				Message:    fmt.Sprintf("failed to create vips image: %s", err),
				StatusCode: http.StatusInternalServerError,
			}
		}
		return res, nil
	default:
		return nil, &Error{
			Message:    "action not allowed",
			StatusCode: http.StatusMethodNotAllowed,
		}
	}
}
