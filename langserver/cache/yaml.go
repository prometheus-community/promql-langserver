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

	d.mu.Lock()
	defer d.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		d.yamlTree = &yamlTree
		return nil
	}
}

func (d *Document) scanYamlTree(ctx context.Context) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return d.scanYamlTreeRec(ctx, d.yamlTree)
	}
}

func (d *Document) scanYamlTreeRec(ctx context.Context, node *yaml.Node) error { //nolint: unparam
	if node == nil {
		return nil
	}

	// Visit all childs
	for _, child := range node.Content {
		err := d.scanYamlTreeRec(ctx, child)
		if err != nil {
			return err
		}
	}

	if node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i < len(node.Content); i += 2 {
		label := node.Content[i]
		value := node.Content[i+1]

		if label.Kind != yaml.ScalarNode || label.Value != "expr" || label.Tag != "!!str" {
			continue
		}

		if value.Kind != yaml.ScalarNode || value.Tag != "!!str" {
			continue
		}

		fmt.Fprintln(os.Stderr, "Found Query:", value)
	}

	return nil
}
