package renderer

import "context"

type Renderer interface {
	// Render draws the diagram and returns the encoded image.
	Render(ctx context.Context, diagram string) ([]byte, error)

	// IsTransient reports whether a failed Render is worth repeating. The
	// implementation owns the meaning of its own failures, the caller owns
	// the retry policy.
	IsTransient(err error) bool
}
