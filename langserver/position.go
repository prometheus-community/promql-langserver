package langserver

import (
	"go/token"

	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
	"github.com/slrtbtfs/go-tools-vendored/span"
)

func (doc *document) positionToProtocolPostion(version float64, pos token.Position) (protocol.Position, bool) {
	doc.Mu.RLock()
	defer doc.Mu.RUnlock()
	if doc.doc.Version > version {
		return protocol.Position{}, false
	}
	line := pos.Line
	char := pos.Column
	// Convert to the Postions as described in the LSP Spec
	point := span.NewPoint(line, char, int(doc.posData.LineStart(line)))
	var err error
	char, err = span.ToUTF16Column(point, []byte(doc.doc.Text))
	char = char - 1
	if err != nil {
		return protocol.Position{}, false
	}
	line = line - 1
	return protocol.Position{
		Line:      float64(line),
		Character: float64(char),
	}, true
}
func (doc *document) protocolPositionToTokenPos(pos protocol.Position) (token.Pos, error) {
	doc.Mu.RLock()
	defer doc.Mu.RUnlock()
	// protocol.Position is 0 based
	line := int(pos.Line) + 1
	char := int(pos.Character)
	point := span.NewPoint(line, 1, int(doc.posData.LineStart(line)))
	point, err := span.FromUTF16Column(point, char, []byte(doc.doc.Text))
	if err != nil {
		return token.NoPos, err
	}
	char = point.Column()
	return doc.posData.LineStart(line) + token.Pos(char), nil

}

func endOfLine(p protocol.Position) protocol.Position {
	return protocol.Position{
		Line:      p.Line + 1,
		Character: 0,
	}
}
