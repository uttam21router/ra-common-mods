package logging

import (
	"context"

	"github.com/routerarchitects/ra-common-mods/ctxmeta"
)

func extractContextFields(ctx context.Context) Fields {
	return ctxmeta.Fields(ctx)
}
