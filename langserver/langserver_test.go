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

package langserver

import (
	"context"
	"testing"

	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/jsonrpc2"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

// TestNotImplemented checks whether unimplemented functions return the approbiate Error
func TestNotImplemented(*testing.T) { // nolint: gocognit, funlen, gocyclo
	s := &server{}

	err := s.DidChangeWorkspaceFolders(context.Background(), &protocol.DidChangeWorkspaceFoldersParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.DidChangeConfiguration(context.Background(), &protocol.DidChangeConfigurationParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.DidSave(context.Background(), &protocol.DidSaveTextDocumentParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.WillSave(context.Background(), &protocol.WillSaveTextDocumentParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.DidChangeWatchedFiles(context.Background(), &protocol.DidChangeWatchedFilesParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.Progress(context.Background(), &protocol.ProgressParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.SelectionRange(context.Background(), &protocol.SelectionRangeParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.SetTraceNotification(context.Background(), &protocol.SetTraceParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	err = s.LogTraceNotification(context.Background(), &protocol.LogTraceParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Implementation(context.Background(), &protocol.ImplementationParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.TypeDefinition(context.Background(), &protocol.TypeDefinitionParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.DocumentColor(context.Background(), &protocol.DocumentColorParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.ColorPresentation(context.Background(), &protocol.ColorPresentationParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.FoldingRange(context.Background(), &protocol.FoldingRangeParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Declaration(context.Background(), &protocol.DeclarationParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.WillSaveWaitUntil(context.Background(), &protocol.WillSaveTextDocumentParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Resolve(context.Background(), &protocol.CompletionItem{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.SignatureHelp(context.Background(), &protocol.SignatureHelpParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Definition(context.Background(), &protocol.DefinitionParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.References(context.Background(), &protocol.ReferenceParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.DocumentHighlight(context.Background(), &protocol.DocumentHighlightParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.DocumentSymbol(context.Background(), &protocol.DocumentSymbolParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.CodeAction(context.Background(), &protocol.CodeActionParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Symbol(context.Background(), &protocol.WorkspaceSymbolParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.CodeLens(context.Background(), &protocol.CodeLensParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.ResolveCodeLens(context.Background(), &protocol.CodeLens{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Formatting(context.Background(), &protocol.DocumentFormattingParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.RangeFormatting(context.Background(), &protocol.DocumentRangeFormattingParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.OnTypeFormatting(context.Background(), &protocol.DocumentOnTypeFormattingParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.Rename(context.Background(), &protocol.RenameParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.PrepareRename(context.Background(), &protocol.PrepareRenameParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.DocumentLink(context.Background(), &protocol.DocumentLinkParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.ResolveDocumentLink(context.Background(), &protocol.DocumentLink{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}

	_, err = s.ExecuteCommand(context.Background(), &protocol.ExecuteCommandParams{})
	if err != nil && err.(*jsonrpc2.Error).Code != jsonrpc2.CodeMethodNotFound {
		panic("Expected a jsonrpc2 Error with CodeMethodNotFound")
	}
}
