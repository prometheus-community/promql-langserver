package b

import (
	_ "github.com/prometheus-community/promql-langserver/vendored/go-tools/lsp/circular/double/one" //@diag("_ \"golang.org/x/tools/internal/lsp/circular/double/one\"", "go list", "import cycle not allowed")
)
