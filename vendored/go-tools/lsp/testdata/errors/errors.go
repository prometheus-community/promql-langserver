package errors

import (
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/types"
)

func _() {
	bob.Bob() //@complete(".")
	types.b //@complete(" //", Bob_interface)
}
