package cache

import (
	"fmt"
	"go/token"
	"os"

	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/span"
)

// e.g. in LineStart
func (d *Document) PositionToProtocolPostion(version float64, pos token.Position) (protocol.Position, bool) {
	d.Mu.RLock()
	defer d.Mu.RUnlock()

	if d.doc.Version > version {
		return protocol.Position{}, false
	}

	line := pos.Line
	char := pos.Column

	// Can happen when parsing empty files
	if line < 1 {
		return protocol.Position{
			Line:      0,
			Character: 0,
		}, true
	}

	// Convert to the Postions as described in the LSP Spec
	// LineStart can panic
	offset := int(d.PosData.LineStart(line)) - d.PosData.Base() + char - 1
	point := span.NewPoint(line, char, offset)

	var err error

	char, err = span.ToUTF16Column(point, []byte(d.doc.Text))
	// Protocol has zero based positions
	char--
	line--

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return protocol.Position{}, false
	}

	return protocol.Position{
		Line:      float64(line),
		Character: float64(char),
	}, true
}

func (d *Document) ProtocolPositionToTokenPos(pos protocol.Position) (token.Pos, error) {
	d.Mu.RLock()
	defer d.Mu.RUnlock()
	// protocol.Position is 0 based
	line := int(pos.Line) + 1
	char := int(pos.Character)
	offset := int(d.PosData.LineStart(line)) - d.PosData.Base()
	point := span.NewPoint(line, 1, offset)
	point, err := span.FromUTF16Column(point, char, []byte(d.doc.Text))

	if err != nil {
		return token.NoPos, err
	}

	char = point.Column()

	return d.PosData.LineStart(line) + token.Pos(char), nil
}

func EndOfLine(p protocol.Position) protocol.Position {
	return protocol.Position{
		Line:      p.Line + 1,
		Character: 0,
	}
}
