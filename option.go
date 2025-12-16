package carapace

import (
	"context"
	"io"
)

type Option func(completer Completer)

type Completer interface {
	SetIn(newIn io.Reader)
	SetOutput(output io.Writer)
	SetErr(newErr io.Writer)

	Context() context.Context
	ExecuteContext(ctx context.Context) error
}
