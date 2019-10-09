package main

import (
	"context"

	"github.com/slrtbtfs/promql-lsp/langserver"
)

func main() {
	ctx, s := langserver.StdioServer(context.Background())
	s.Run(ctx)
}
