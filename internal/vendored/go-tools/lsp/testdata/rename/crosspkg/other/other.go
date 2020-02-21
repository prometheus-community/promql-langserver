package other

import "github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/rename/crosspkg"

func Other() {
	crosspkg.Bar
	crosspkg.Foo()
}
