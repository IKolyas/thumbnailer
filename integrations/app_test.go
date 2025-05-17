package integrations

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL     = "http://localhost:8080"
	testTimeout = 5 * time.Second
)

func TestAPISuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	sourceHost := os.Getenv("SRC_HOST")
	url := buildURL("/fill/500/500/", sourceHost, "/buket.jpg")

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err, "failed to create request")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "request failed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "image/jpeg", resp.Header.Get("Content-Type"))
}

func TestAPINotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	sourceHost := os.Getenv("SRC_HOST")
	url := buildURL("/fill/600/600/", sourceHost, "/not-found.jpg")

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAPINoImage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	sourceHost := os.Getenv("SRC_HOST")
	url := buildURL("/fill/600/600/", sourceHost, "/no-image.txt")

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
}

func TestAPIErrorServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	sourceHost := os.Getenv("SRC_HOST")
	url := buildURL("/fill/600/600/", sourceHost, "/error.txt")

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func buildURL(paths ...string) string {
	return baseURL + joinPaths(paths...)
}

func joinPaths(paths ...string) string {
	var result string
	for _, path := range paths {
		result += path
	}
	return result
}
