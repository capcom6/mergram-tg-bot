package mermaidinkfx_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/capcom6/mergram-tg-bot/pkg/mermaidinkfx"
)

const testImage = "image-bytes"

func newTestClient(t *testing.T, baseURL string) *mermaidinkfx.Client {
	t.Helper()

	client, err := mermaidinkfx.NewClient(
		mermaidinkfx.Config{BaseURL: baseURL},
		http.DefaultClient,
	)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	return client
}

func TestRenderReturnsImage(t *testing.T) {
	t.Parallel()

	var requests atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)

		if !strings.HasPrefix(r.URL.Path, "/img/pako:") {
			t.Errorf("unexpected path %q", r.URL.Path)
		}

		_, _ = w.Write([]byte(testImage))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)

	data, err := client.Render(t.Context(), "graph TD\nA-->B")
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if string(data) != testImage {
		t.Errorf("expected %q, got %q", testImage, string(data))
	}

	if got := requests.Load(); got != 1 {
		t.Errorf("expected 1 request, got %d", got)
	}
}

func TestRenderReportsPermanentFailure(t *testing.T) {
	t.Parallel()

	var requests atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)

	_, err := client.Render(t.Context(), "not a diagram")
	if err == nil {
		t.Fatal("expected an error")
	}

	if !errors.Is(err, mermaidinkfx.ErrAPIError) {
		t.Errorf("expected ErrAPIError, got %v", err)
	}

	if client.IsTransient(err) {
		t.Errorf("expected a permanent error, got %v", err)
	}

	if got := requests.Load(); got != 1 {
		t.Errorf("expected 1 request, got %d", got)
	}
}

func TestRenderReportsTransientFailure(t *testing.T) {
	t.Parallel()

	var requests atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)

	_, err := client.Render(t.Context(), "graph TD\nA-->B")
	if err == nil {
		t.Fatal("expected an error")
	}

	if !errors.Is(err, mermaidinkfx.ErrUpstreamUnavailable) {
		t.Errorf("expected ErrUpstreamUnavailable, got %v", err)
	}

	if !client.IsTransient(err) {
		t.Errorf("expected a transient error, got %v", err)
	}

	if got := requests.Load(); got != 1 {
		t.Errorf("expected 1 request, got %d", got)
	}
}

func TestNewClientRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	_, err := mermaidinkfx.NewClient(mermaidinkfx.Config{BaseURL: ""}, http.DefaultClient)
	if err == nil {
		t.Fatal("expected an error")
	}

	if !errors.Is(err, mermaidinkfx.ErrInvalidConfig) {
		t.Errorf("expected ErrInvalidConfig, got %v", err)
	}
}
