package unimported

import (
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/baz"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/signature" // provide type information for unimported completions in the other file
)

func _() {
	foo.StructFoo{} //@item(litFooStructFoo, "foo.StructFoo{}", "struct{...}", "struct")

	// We get the literal completion for "foo.StructFoo{}" even though we haven't
	// imported "foo" yet.
	baz.FooStruct = f //@snippet(" //", litFooStructFoo, "foo.StructFoo{$0\\}", "foo.StructFoo{$0\\}")
}
