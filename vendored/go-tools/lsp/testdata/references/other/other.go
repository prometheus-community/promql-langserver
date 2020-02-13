package other

import (
	references "github.com/prometheus-community/promql-langserver/vendored/go-tools/lsp/references"
)

func _() {
	references.Q = "hello" //@mark(assignExpQ, "Q")
	bob := func(_ string) {}
	bob(references.Q) //@mark(bobExpQ, "Q")
}
