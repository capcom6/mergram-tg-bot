package renderer

import "context"

type Renderer interface {
	Render(ctx context.Context, diagram string) ([]byte, error)
}
