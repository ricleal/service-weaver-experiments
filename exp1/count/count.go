package count

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

// counter must implement the Counter interface.
var _ Counter = (*counter)(nil)

// Counter component.
type Counter interface {
	Count(context.Context, string) (int, error)
}

// Implementation of the Counter component.
type counter struct {
	weaver.Implements[Counter]
}

func (c *counter) Count(_ context.Context, s string) (int, error) {
	return len(s), nil
}
