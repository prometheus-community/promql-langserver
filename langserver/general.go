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
	"os"

	"github.com/slrtbtfs/go-tools-vendored/jsonrpc2"
	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
)

func (s *Server) Initialize(ctx context.Context, params *protocol.ParamInitia) (*protocol.InitializeResult, error) {
	s.stateMu.Lock()
	state := s.state
	s.stateMu.Unlock()
	if state >= serverInitializing {
		return nil, jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidRequest, "server already initialized")
	}
	s.stateMu.Lock()
	s.state = serverInitializing
	s.stateMu.Unlock()

	s.cache.init()

	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: true,
				// Support incremental changes
				Change: 2, 
			},
			HoverProvider:      false,
			CompletionProvider: nil,
			SignatureHelpProvider: &protocol.SignatureHelpOptions{
				TriggerCharacters: nil,
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
					WorkDoneProgress: false,
				},
			},
			DefinitionProvider:              false,
			ReferencesProvider:              false,
			DocumentHighlightProvider:       false,
			DocumentSymbolProvider:          false,
			CodeActionProvider:              nil,
			WorkspaceSymbolProvider:         false,
			CodeLensProvider:                nil,
			DocumentFormattingProvider:      false,
			DocumentRangeFormattingProvider: false,
			DocumentOnTypeFormattingProvider: &protocol.DocumentOnTypeFormattingOptions{
				FirstTriggerCharacter: "",
				MoreTriggerCharacter:  nil,
			},
			RenameProvider:       nil,
			DocumentLinkProvider: nil,
			ExecuteCommandProvider: &protocol.ExecuteCommandOptions{
				Commands: nil,
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
					WorkDoneProgress: false,
				},
			},
			Experimental:           nil,
			ImplementationProvider: false,
			TypeDefinitionProvider: false,
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

func (s *Server) Initialized(ctx context.Context, params *protocol.InitializedParams) error {
	s.stateMu.Lock()
	s.state = serverInitialized
	s.stateMu.Unlock()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	if s.state < serverInitialized {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidRequest, "server not initialized")
	}

	s.state = serverShutDown
	return nil
}
func (s *Server) Exit(ctx context.Context) error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	if s.state != serverShutDown {
		os.Exit(1)
	}
	os.Exit(0)
	return nil
}
