package integrations

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL     = "http://localhost:8080/fill/500/500/storage-server"
	testTimeout = 5 * time.Second
)

func TestAPISuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/buket.jpg", nil)
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

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/not-found.jpg", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAPINoImage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/no-image.txt", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
}

func TestAPIErrorServer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/error.txt", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
