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

	"github.com/slrtbtfs/prometheus/promql"
)

// CompiledQuery stores the results of compiling one query
type CompiledQuery struct {
	Ast     promql.Node
	Err     *promql.ParseErr
	Content string
	Record  string
}

func (d *Document) compile(ctx context.Context) error {
	defer d.compilers.Done()

	switch d.GetLanguageID() {
	case "promql":
		d.compilers.Add(1)
		return d.compileQuery(ctx, true, 0, 0, "")
	case "yaml":
		err := d.parseYamls(ctx)
		if err != nil {
			return err
		}

		d.compilers.Add(1)

		err = d.scanYamlTree(ctx)
		if err != nil {
			return err
		}
	default:
		fmt.Fprintf(os.Stderr, "Unsupported Filetype: %s\n", d.GetLanguageID())
	}

	return nil
}

// compileQuery compiles the query at the position given by the last two arguments
// if fullFile is set, the last two arguments are ignored and the full file is assumed
// to be one query
func (d *Document) compileQuery(ctx context.Context, fullFile bool, pos token.Pos, endPos token.Pos, record string) error { //nolint:lll
	defer d.compilers.Done()

	var content string

	var expired error

	if fullFile {
		content, expired = d.GetContent(ctx)
		pos = token.Pos(d.posData.Base())
	} else {
		content, expired = d.GetSubstring(ctx, pos, endPos)
	}

	if expired != nil {
		return expired
	}

	file := d.posData

	ast, err := promql.ParsePartOfFile(content, file, pos)

	var parseErr *promql.ParseErr

	var ok bool

	if parseErr, ok = err.(*promql.ParseErr); !ok {
		parseErr = nil
	}

	err = d.AddCompileResult(ctx, ast, parseErr, record, content)
	if err != nil {
		return err
	}

	if parseErr != nil {
		diagnostic, err := d.promQLErrToProtocolDiagnostic(ctx, parseErr)
		if err != nil {
			return err
		}

		err = d.AddDiagnostic(ctx, diagnostic)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddCompileResult updates the compilation Results of a Document. Discards the Result if the context is expired
func (d *Document) AddCompileResult(ctx context.Context, ast promql.Node, err *promql.ParseErr, record string, content string) error { //nolint: lll
	d.mu.Lock()
	defer d.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		d.queries = append(d.queries, &CompiledQuery{ast, err, content, record})
		return nil
	}
}
