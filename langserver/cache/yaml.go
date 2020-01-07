// Copyright 2019 Tobias Guggenmos
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"errors"
	"go/token"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

// YamlDoc contains the results of compiling a yaml document
type YamlDoc struct {
	AST yaml.Node
	Err error
	// Not encoded in the AST
	End token.Pos
	// Offset that has to be added to every line number before translating into a token.Pos
	LineOffset int
}

func (d *DocumentHandle) parseYamls() error {
	content, err := d.GetContent()
	if err != nil {
		return err
	}

	reader := strings.NewReader(content)

	lineOffset := 0

	for unread := reader.Len(); unread > 0; {
		var yamlDoc YamlDoc

		decoder := yaml.NewDecoder(reader)

		yamlDoc.Err = decoder.Decode(&yamlDoc.AST)

		unread = reader.Len()

		yamlDoc.End = token.Pos(d.doc.posData.Base() + len(content) - unread)
		yamlDoc.LineOffset = lineOffset

		// Update Line Offset for the next document
		lineOffset = d.doc.posData.Line(yamlDoc.End) - 1

		err := d.addYaml(&yamlDoc)
		if err != nil {
			return err
		}

		if errors.Is(yamlDoc.Err, io.EOF) {
			return yamlDoc.Err
		}
	}

	return nil
}

func (d *DocumentHandle) addYaml(yaml *YamlDoc) error {
	d.doc.mu.Lock()
	defer d.doc.mu.Unlock()

	select {
	case <-d.ctx.Done():
		return d.ctx.Err()
	default:
		d.doc.yamls = append(d.doc.yamls, yaml)

		return nil
	}
}

func (d *DocumentHandle) scanYamlTree() error {
	defer d.doc.compilers.Done()

	yamls, err := d.GetYamls()
	if err != nil {
		return err
	}

	for _, yamlDoc := range yamls {
		err := d.scanYamlTreeRec(&yamlDoc.AST, yamlDoc.End, yamlDoc.LineOffset, nil)
		if err != nil {
			return err
		}
	}

	return err
}

// nolint
func (d *DocumentHandle) scanYamlTreeRec(node *yaml.Node, nodeEnd token.Pos, lineOffset int, path []string) error { //nolint: unparam
	if node == nil {
		return nil
	}

	// Visit all childs
	for i, child := range node.Content {
		var err error

		var childEnd token.Pos

		var childPath []string

		if i+1 < len(node.Content) && node.Content[i+1] != nil {
			next := node.Content[i+1]

			childEnd, err = d.YamlPositionToTokenPos(next.Line, next.Column, lineOffset)
			if err != nil {
				return err
			}
		} else {
			childEnd = nodeEnd
		}

		if node.Value != "" {
			childPath = append(childPath, node.Value)
		}

		if node.Kind == yaml.MappingNode && i > 0 && i%2 == 1 {
			childPath = append(childPath, node.Content[i-1].Value)
		}

		err = d.scanYamlTreeRec(child, childEnd, lineOffset, append(path, childPath...))
		if err != nil {
			return err
		}
	}

	if relevantYamlPath(path) {
		if err := d.foundRelevantYamlPath(node, nodeEnd, lineOffset); err != nil {
			return err
		}
	}
	return nil
}

//nolint:gocognit
func (d *DocumentHandle) foundRelevantYamlPath(node *yaml.Node, nodeEnd token.Pos, lineOffset int) error { //nolint: lll
	if node.Kind != yaml.MappingNode {
		return nil
	}

	var expr *yaml.Node

	var exprEnd token.Pos

	var record *yaml.Node

	for i := 0; i+1 < len(node.Content); i += 2 {
		label := node.Content[i]
		value := node.Content[i+1]

		if label == nil || label.Kind != yaml.ScalarNode || label.Tag != "!!str" {
			continue
		}

		if value == nil || value.Kind != yaml.ScalarNode || value.Tag != "!!str" {
			continue
		}

		switch label.Value {
		case "expr":
			var err error

			if i+2 < len(node.Content) && node.Content[i+2] != nil {
				next := node.Content[i+2]

				exprEnd, err = d.YamlPositionToTokenPos(next.Line, next.Column, lineOffset)
				if err != nil {
					return err
				}
			} else {
				exprEnd = nodeEnd
			}

			expr = value
		case "record":
			record = value
		}
	}

	if expr == nil {
		return nil
	}

	err := d.foundQuery(expr, exprEnd, record, lineOffset)

	if err != nil {
		return err
	}

	return nil
}

func relevantYamlPath(path []string) bool {
	relevantSuffixes := [][]string{
		{"alerts"},
		{"groups", "rules"},
		{"recordingrule"},
	}

OUTER:
	for _, suffix := range relevantSuffixes {
		if len(suffix) > len(path) {
			continue
		}

		shortPath := path[len(path)-len(suffix):]

		for i := range suffix {
			if suffix[i] != shortPath[i] {
				continue OUTER
			}
		}

		return true
	}

	return false
}

func (d *DocumentHandle) foundQuery(node *yaml.Node, endPos token.Pos, record *yaml.Node, lineOffset int) error { //nolint: lll
	line := node.Line
	col := node.Column

	if node.Style == yaml.LiteralStyle || node.Style == yaml.FoldedStyle {
		// The query starts on the line following the '|' or '>'
		line++

		col = 1
	}

	pos, err := d.YamlPositionToTokenPos(line, col, lineOffset)
	if err != nil {
		return err
	}

	if node.Style == yaml.SingleQuotedStyle || node.Style == yaml.DoubleQuotedStyle {
		err = d.warnQuotedYaml(pos, endPos)
		return err
	}

	d.doc.compilers.Add(1)

	var recordValue string

	if record != nil {
		recordValue = record.Value
	}

	go d.compileQuery(false, pos, endPos, recordValue) //nolint: errcheck

	return nil
}
