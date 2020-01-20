package other

import "github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/rename/crosspkg"

func Other() {
	crosspkg.Bar
	crosspkg.Foo()
}
