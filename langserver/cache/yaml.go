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
	"context"
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

func (d *Document) parseYamls(ctx context.Context) error {
	content, err := d.GetContent(ctx)
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

		yamlDoc.End = token.Pos(d.posData.Base() + len(content) - unread)
		yamlDoc.LineOffset = lineOffset

		// Update Line Offset for the next document
		lineOffset = d.posData.Line(yamlDoc.End) - 1

		err := d.addYaml(ctx, &yamlDoc)
		if err != nil {
			return err
		}

		if errors.Is(yamlDoc.Err, io.EOF) {
			return yamlDoc.Err
		}
	}

	return nil
}

func (d *Document) addYaml(ctx context.Context, yaml *YamlDoc) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		d.yamls = append(d.yamls, yaml)

		return nil
	}
}

func (d *Document) scanYamlTree(ctx context.Context) error {
	defer d.compilers.Done()

	yamls, err := d.GetYamls(ctx)
	if err != nil {
		return err
	}

	for _, yamlDoc := range yamls {
		err := d.scanYamlTreeRec(ctx, &yamlDoc.AST, yamlDoc.End, yamlDoc.LineOffset, nil)
		if err != nil {
			return err
		}
	}

	return err
}

// nolint
func (d *Document) scanYamlTreeRec(ctx context.Context, node *yaml.Node, nodeEnd token.Pos, lineOffset int, path []string) error { //nolint: unparam
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

			childEnd, err = d.YamlPositionToTokenPos(ctx, next.Line, next.Column, lineOffset)
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

		err = d.scanYamlTreeRec(ctx, child, childEnd, lineOffset, append(path, childPath...))
		if err != nil {
			return err
		}
	}

	if node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i+1 < len(node.Content); i += 2 {
		label := node.Content[i]
		value := node.Content[i+1]

		if label == nil || label.Kind != yaml.ScalarNode || label.Value != "expr" || label.Tag != "!!str" {
			continue
		}

		if value == nil || value.Kind != yaml.ScalarNode || value.Tag != "!!str" {
			continue
		}

		var err error

		var valueEnd token.Pos

		if i+2 < len(node.Content) && node.Content[i+2] != nil {
			next := node.Content[i+2]

			valueEnd, err = d.YamlPositionToTokenPos(ctx, next.Line, next.Column, lineOffset)
			if err != nil {
				return err
			}
		} else {
			valueEnd = nodeEnd
		}

		err = d.foundQuery(ctx, value, valueEnd, lineOffset)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Document) foundQuery(ctx context.Context, node *yaml.Node, endPos token.Pos, lineOffset int) error {
	line := node.Line
	col := node.Column

	if node.Style == yaml.LiteralStyle || node.Style == yaml.FoldedStyle {
		// The query starts on the line following the '|' or '>'
		line++

		col = 1
	}

	pos, err := d.YamlPositionToTokenPos(ctx, line, col, lineOffset)
	if err != nil {
		return err
	}

	if node.Style == yaml.SingleQuotedStyle || node.Style == yaml.DoubleQuotedStyle {
		err = d.warnQuotedYaml(ctx, pos, endPos)
		return err
	}

	d.compilers.Add(1)

	go d.compileQuery(ctx, false, pos, endPos) //nolint: errcheck

	return nil
}
