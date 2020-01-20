package other

import (
	references "github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/references"
)

func _() {
	references.Q = "hello" //@mark(assignExpQ, "Q")
	bob := func(_ string) {}
	bob(references.Q) //@mark(bobExpQ, "Q")
}
