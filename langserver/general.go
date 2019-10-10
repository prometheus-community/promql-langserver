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
			},
			HoverProvider: false,
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters:   nil,
				AllCommitCharacters: nil,
				ResolveProvider:     false,
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
					WorkDoneProgress: false,
				},
			},
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
