package a

import (
	_ "github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/circular/triple/b" //@diag("_ \"golang.org/x/tools/internal/lsp/circular/triple/b\"", "go list", "import cycle not allowed")
)
