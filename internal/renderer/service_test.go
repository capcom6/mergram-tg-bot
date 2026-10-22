package renderer_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/alitto/pond/v2"
	"github.com/capcom6/mergram-tg-bot/internal/renderer"
	"go.uber.org/zap"
)

const (
	testMaxAttempts = 3
	testBackoff     = 10 * time.Millisecond
	testTimeout     = 2 * time.Second
	testImage       = "image-bytes"
)

// errTransient is the failure the stub renderer reports as worth retrying,
// errPermanent the one it does not.
var (
	errTransient = errors.New("transient failure")
	errPermanent = errors.New("permanent failure")
)

// stubRenderer is a renderer.Renderer that fails with the given errors, one
// per call, and returns the image once they are exhausted.
type stubRenderer struct {
	mu    sync.Mutex
	errs  []error
	calls int
}

func (s *stubRenderer) Render(_ context.Context, _ string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.calls++

	if len(s.errs) == 0 {
		return []byte(testImage), nil
	}

	err := s.errs[0]
	s.errs = s.errs[1:]

	return nil, err
}

func (s *stubRenderer) IsTransient(err error) bool {
	return errors.Is(err, errTransient)
}

func (s *stubRenderer) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.calls
}

func newService(t *testing.T, stub renderer.Renderer) *renderer.Service {
	t.Helper()

	svc, err := renderer.NewService(
		renderer.Config{
			Timeout:     testTimeout,
			MaxAttempts: testMaxAttempts,
			Backoff:     testBackoff,
		},
		stub,
		pond.NewPool(1),
		zap.NewNop(),
	)
	if err != nil {
		t.Fatalf("create service: %v", err)
	}

	return svc
}

func TestRenderSucceedsOnFirstAttempt(t *testing.T) {
	t.Parallel()

	stub := &stubRenderer{}
	svc := newService(t, stub)

	data, err := svc.Render(t.Context(), "graph TD\nA-->B")
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if string(data) != testImage {
		t.Errorf("expected %q, got %q", testImage, string(data))
	}

	if got := stub.callCount(); got != 1 {
		t.Errorf("expected 1 call, got %d", got)
	}
}

func TestRenderRetriesTransientFailure(t *testing.T) {
	t.Parallel()

	stub := &stubRenderer{errs: []error{errTransient}}
	svc := newService(t, stub)

	data, err := svc.Render(t.Context(), "graph TD\nA-->B")
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if string(data) != testImage {
		t.Errorf("expected %q, got %q", testImage, string(data))
	}

	if got := stub.callCount(); got != 2 {
		t.Errorf("expected 2 calls, got %d", got)
	}
}

func TestRenderDoesNotRetryPermanentFailure(t *testing.T) {
	t.Parallel()

	stub := &stubRenderer{errs: []error{errPermanent}}
	svc := newService(t, stub)

	_, err := svc.Render(t.Context(), "not a diagram")
	if err == nil {
		t.Fatal("expected an error")
	}

	if !errors.Is(err, errPermanent) {
		t.Errorf("expected ErrAPIError, got %v", err)
	}

	if got := stub.callCount(); got != 1 {
		t.Errorf("expected 1 call, got %d", got)
	}
}

func TestRenderGivesUpAfterMaxAttempts(t *testing.T) {
	t.Parallel()

	stub := &stubRenderer{errs: []error{
		errTransient,
		errTransient,
		errTransient,
	}}
	svc := newService(t, stub)

	_, err := svc.Render(t.Context(), "graph TD\nA-->B")
	if err == nil {
		t.Fatal("expected an error")
	}

	if !errors.Is(err, errTransient) {
		t.Errorf("expected ErrUpstreamUnavailable, got %v", err)
	}

	if got := stub.callCount(); got != testMaxAttempts {
		t.Errorf("expected %d calls, got %d", testMaxAttempts, got)
	}
}

func TestRenderRetriesAfterAttemptTimeout(t *testing.T) {
	t.Parallel()

	// The first attempt hangs until its own budget is over, the second one
	// answers right away.
	stub := &hangingRenderer{}
	svc := newService(t, stub)

	data, err := svc.Render(t.Context(), "graph TD\nA-->B")
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if string(data) != testImage {
		t.Errorf("expected %q, got %q", testImage, string(data))
	}

	if got := stub.callCount(); got != 2 {
		t.Errorf("expected 2 calls, got %d", got)
	}
}

func TestRenderStopsWhenContextIsCanceled(t *testing.T) {
	t.Parallel()

	stub := &stubRenderer{errs: []error{errTransient}}
	svc := newService(t, stub)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := svc.Render(ctx, "graph TD\nA-->B")
	if err == nil {
		t.Fatal("expected an error")
	}

	if got := stub.callCount(); got > 1 {
		t.Errorf("expected no retry, got %d calls", got)
	}
}

func TestIsTransient(t *testing.T) {
	t.Parallel()

	stub := &stubRenderer{}
	svc := newService(t, stub)

	if !svc.IsTransient(errTransient) {
		t.Error("expected a transient failure to be transient")
	}

	if svc.IsTransient(errPermanent) {
		t.Error("expected a permanent failure not to be transient")
	}
}

func TestNewServiceRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	_, err := renderer.NewService(
		renderer.Config{Timeout: testTimeout, MaxAttempts: 0, Backoff: testBackoff},
		&stubRenderer{},
		pond.NewPool(1),
		zap.NewNop(),
	)
	if err == nil {
		t.Fatal("expected an error")
	}

	if !errors.Is(err, renderer.ErrInvalidConfig) {
		t.Errorf("expected ErrInvalidConfig, got %v", err)
	}
}

// hangingRenderer blocks the first call until its context is done and answers
// every later call.
type hangingRenderer struct {
	mu    sync.Mutex
	calls int
}

func (h *hangingRenderer) Render(ctx context.Context, _ string) ([]byte, error) {
	h.mu.Lock()
	h.calls++
	first := h.calls == 1
	h.mu.Unlock()

	if first {
		<-ctx.Done()
		return nil, ctx.Err()
	}

	return []byte(testImage), nil
}

func (h *hangingRenderer) IsTransient(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
}

func (h *hangingRenderer) callCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.calls
}
