package other

import "github.com/prometheus-community/promql-langserver/vendored/go-tools/lsp/rename/crosspkg"

func Other() {
	crosspkg.Bar
	crosspkg.Foo()
}
