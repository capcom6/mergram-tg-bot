package pondfx

import "github.com/alitto/pond/v2"

func New(config Config, options []pond.Option) pond.Pool {
	return pond.NewPool(config.MaxConcurrency, options...)
}
