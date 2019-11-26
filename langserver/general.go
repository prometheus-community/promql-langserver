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
	"errors"

	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/jsonrpc2"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

// Initialize handles a call from the client to initialize the server
// required by the protocol.Server interface
// nolint:funlen
func (s *server) Initialize(ctx context.Context, params *protocol.ParamInitia) (*protocol.InitializeResult, error) {
	s.stateMu.Lock()
	state := s.state
	s.stateMu.Unlock()

	if state != serverCreated {
		return nil, jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidRequest, "server already initialized")
	}

	s.stateMu.Lock()
	s.state = serverInitializing
	s.stateMu.Unlock()

	s.cache.Init()

	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: true,
				// Support incremental changes
				Change: 2,
			},
			HoverProvider: true,
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{
					" ", "\n", "\t", "(", ")", "[", "]", "{", "}", "+", "-", "*", "/", "!", "=", "\"", ",",
				},
				AllCommitCharacters: nil,
				ResolveProvider:     false,
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
					WorkDoneProgress: false,
				},
			},
			SignatureHelpProvider:            nil,
			DefinitionProvider:               false,
			ReferencesProvider:               false,
			DocumentHighlightProvider:        false,
			DocumentSymbolProvider:           false,
			CodeActionProvider:               nil,
			WorkspaceSymbolProvider:          false,
			CodeLensProvider:                 nil,
			DocumentFormattingProvider:       false,
			DocumentRangeFormattingProvider:  false,
			DocumentOnTypeFormattingProvider: nil,
			RenameProvider:                   nil,
			DocumentLinkProvider:             nil,
			ExecuteCommandProvider:           nil,
			Experimental:                     nil,
			ImplementationProvider:           false,
			TypeDefinitionProvider:           false,
			Workspace: &struct {
				WorkspaceFolders *struct {
					Supported           bool   "json:\"supported,omitempty\""
					ChangeNotifications string "json:\"changeNotifications,omitempty\""
				} "json:\"workspaceFolders,omitempty\""
			}{
				WorkspaceFolders: &struct {
					Supported           bool   "json:\"supported,omitempty\""
					ChangeNotifications string "json:\"changeNotifications,omitempty\""
				}{
					Supported:           false,
					ChangeNotifications: "",
				},
			},
			ColorProvider:          false,
			FoldingRangeProvider:   false,
			DeclarationProvider:    false,
			SelectionRangeProvider: false,
		},
	}, nil
}

// Initialized receives a confirmation by the client that the connection has been initialized
// required by the protocol.Server interface
func (s *server) Initialized(ctx context.Context, params *protocol.InitializedParams) error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	if s.state != serverInitializing {
		return errors.New("cannot initialize server: wrong server state")
	}

	s.state = serverInitialized

	return nil
}

// Shutdown receives a call from the client to shutdown the connection
// required by the protocol.Server interface
func (s *server) Shutdown(ctx context.Context) error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	if s.state != serverInitialized {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidRequest, "server not initialized")
	}

	s.state = serverShutDown

	return nil
}

// Exit ends the connection
// required by the protocol.Server interface
func (s *server) Exit(ctx context.Context) error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	if s.state != serverShutDown {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidRequest, "server not shutdown")
	}

	s.exit()

	return nil
}
