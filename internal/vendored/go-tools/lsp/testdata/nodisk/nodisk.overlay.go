package nodisk

import (
	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/foo"
)

func _() {
	foo.Foo() //@complete("F", Foo, IntFoo, StructFoo)
}