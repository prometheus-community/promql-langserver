package a

import (
	_ "github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/circular/triple/b" //@diag("_ \"golang.org/x/tools/internal/lsp/circular/triple/b\"", "go list", "import cycle not allowed")
)
