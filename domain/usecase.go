package domain

import "context"

// Usecase ...
type Usecase[I, O any] interface {
	Execute(ctx context.Context, input I) (output *O, err error)
}

type UsecaseWithoutOutput[I any] interface {
	Execute(ctx context.Context, input I) (err error)
}
