package nodisk

import (
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/foo"
)

func _() {
	foo.Foo() //@complete("F", Foo, IntFoo, StructFoo)
}