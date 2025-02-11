package user

import (
	"context"
	"github/kijunpos/app/lib/apm"
)

func (c *Command) Create(ctx context.Context) (err error) {
	ctx, span := apm.GetTracer().Start(ctx, "repository.user.Create")
	defer span.End()

	return nil
}
