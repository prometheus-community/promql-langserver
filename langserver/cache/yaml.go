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
	"fmt"
	"go/token"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func (d *Document) parseYaml(ctx context.Context) error {
	content, err := d.GetContent(ctx)
	if err != nil {
		return err
	}

	var yamlTree yaml.Node

	reader := strings.NewReader(content)
	decoder := yaml.NewDecoder(reader)

	err = decoder.Decode(&yamlTree)
	if err != nil {
		return err
	}

	unread := reader.Len()
	yamlEnd := d.posData.Base() + len(content) - unread

	d.mu.Lock()
	defer d.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		d.yamlTree = &yamlTree
		d.yamlEnd = token.Pos(yamlEnd)

		return nil
	}
}

func (d *Document) scanYamlTree(ctx context.Context) error {
	defer d.compilers.Done()

	yamlTree, yamlEnd, err := d.GetYamlTree(ctx)
	if err != nil {
		return err
	}

	return d.scanYamlTreeRec(ctx, yamlTree, yamlEnd)
}

// nolint
func (d *Document) scanYamlTreeRec(ctx context.Context, node *yaml.Node, nodeEnd token.Pos) error { //nolint: unparam
	if node == nil {
		return nil
	}

	// Visit all childs
	for i, child := range node.Content {
		var err error

		var childEnd token.Pos

		if i+1 < len(node.Content) && node.Content[i+1] != nil {
			next := node.Content[i+1]

			childEnd, err = d.yamlPositionToTokenPos(ctx, next.Line, next.Column)
			if err != nil {
				return err
			}
		} else {
			childEnd = nodeEnd
		}

		err = d.scanYamlTreeRec(ctx, child, childEnd)
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

			valueEnd, err = d.yamlPositionToTokenPos(ctx, next.Line, next.Column)
			if err != nil {
				return err
			}
		} else {
			valueEnd = nodeEnd
		}

		d.foundQuery(ctx, value, valueEnd)
	}

	return nil
}

func (d *Document) foundQuery(ctx context.Context, node *yaml.Node, endPos token.Pos) error {
	line := node.Line
	col := node.Column

	if node.Style == yaml.LiteralStyle || node.Style == yaml.FoldedStyle {
		// The query starts on the line following the '|' or '>'
		line++

		col = 1
	} else if node.Style == yaml.SingleQuotedStyle || node.Style == yaml.DoubleQuotedStyle {
		fmt.Fprintf(os.Stderr, "Line %d: Warning: Quoted queries are not supported\n", node.Line)
		return nil
	}

	pos, err := d.yamlPositionToTokenPos(ctx, line, col)
	if err != nil {
		return err
	}

	d.compilers.Add(1)

	go d.compileQuery(ctx, false, pos, endPos)

	return nil
}
