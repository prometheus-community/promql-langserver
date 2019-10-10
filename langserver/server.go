package langserver

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/slrtbtfs/go-tools-vendored/jsonrpc2"
	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
)

type Server struct {
	Conn   *jsonrpc2.Conn
	client protocol.Client

	state   serverState
	stateMu sync.Mutex
}

type serverState int

const (
	serverCreated      = serverState(iota)
	serverInitializing // set once the server has received "initialize" request
	serverInitialized  // set once the server has received "initialized" request
	serverShutDown
)

func (s *Server) Run(ctx context.Context) error {
	return s.Conn.Run(ctx)
}

func ServerFromStream(ctx context.Context, stream jsonrpc2.Stream) (context.Context, *Server) {
	s := &Server{}
	ctx, s.Conn, s.client = protocol.NewServer(ctx, stream, s)
	return ctx, s
}

func StdioServer(ctx context.Context) (context.Context, *Server) {
	stream := jsonrpc2.NewHeaderStream(os.Stdin, os.Stdout)
	stream = protocol.LoggingStream(stream, os.Stderr)
	return ServerFromStream(ctx, stream)
}

func notImplemented(method string) *jsonrpc2.Error {
	err := jsonrpc2.NewErrorf(jsonrpc2.CodeMethodNotFound, "method %q no yet implemented", method)
	fmt.Fprint(os.Stderr, err.Error())
	return err
}
func (s *Server) DidChangeWorkspaceFolders(_ context.Context, _ *protocol.DidChangeWorkspaceFoldersParams) error {
	return notImplemented("DidChangeWorkspaceFolders")
}

func (s *Server) DidChangeConfiguration(_ context.Context, _ *protocol.DidChangeConfigurationParams) error {
	return notImplemented("DidChangeConfiguration")
}

func (s *Server) DidOpen(_ context.Context, _ *protocol.DidOpenTextDocumentParams) error {
	return notImplemented("DidOpen")
}

func (s *Server) DidChange(_ context.Context, _ *protocol.DidChangeTextDocumentParams) error {
	return notImplemented("DidChange")
}

func (s *Server) DidClose(_ context.Context, _ *protocol.DidCloseTextDocumentParams) error {
	return notImplemented("DidClose")
}

func (s *Server) DidSave(_ context.Context, _ *protocol.DidSaveTextDocumentParams) error {
	return notImplemented("DidSave")
}

func (s *Server) WillSave(_ context.Context, _ *protocol.WillSaveTextDocumentParams) error {
	return notImplemented("WillSave")
}

func (s *Server) DidChangeWatchedFiles(_ context.Context, _ *protocol.DidChangeWatchedFilesParams) error {
	return notImplemented("DidChangeWatchedFiles")
}

func (s *Server) Progress(_ context.Context, _ *protocol.ProgressParams) error {
	return notImplemented("Progress")
}

func (s *Server) SetTraceNotification(_ context.Context, _ *protocol.SetTraceParams) error {
	return notImplemented("SetTraceNotification")
}

func (s *Server) LogTraceNotification(_ context.Context, _ *protocol.LogTraceParams) error {
	return notImplemented("LogTraceNotification")
}

func (s *Server) Implementation(_ context.Context, _ *protocol.ImplementationParams) ([]protocol.Location, error) {
	return nil, notImplemented("Implementation")
}

func (s *Server) TypeDefinition(_ context.Context, _ *protocol.TypeDefinitionParams) ([]protocol.Location, error) {
	return nil, notImplemented("TypeDefinition")
}

func (s *Server) DocumentColor(_ context.Context, _ *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	return nil, notImplemented("DocumentColor")
}

func (s *Server) ColorPresentation(_ context.Context, _ *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	return nil, notImplemented("ColorPresentation")
}

func (s *Server) FoldingRange(_ context.Context, _ *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	return nil, notImplemented("FoldingRange")
}

func (s *Server) Declaration(_ context.Context, _ *protocol.DeclarationParams) ([]protocol.DeclarationLink, error) {
	return nil, notImplemented("Declaration")
}

func (s *Server) SelectionRange(_ context.Context, _ *protocol.SelectionRangeParams) ([]protocol.SelectionRange, error) {
	return nil, notImplemented("SelectionRange")
}

func (s *Server) WillSaveWaitUntil(_ context.Context, _ *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	return nil, notImplemented("WillSaveWaitUntil")
}

func (s *Server) Completion(_ context.Context, _ *protocol.CompletionParams) (*protocol.CompletionList, error) {
	return nil, notImplemented("Completion")
}

func (s *Server) Resolve(_ context.Context, _ *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return nil, notImplemented("Resolve")
}

func (s *Server) Hover(_ context.Context, _ *protocol.HoverParams) (*protocol.Hover, error) {
	return nil, notImplemented("Hover")
}

func (s *Server) SignatureHelp(_ context.Context, _ *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	return nil, notImplemented("SignatureHelp")
}

func (s *Server) Definition(_ context.Context, _ *protocol.DefinitionParams) ([]protocol.Location, error) {
	return nil, notImplemented("Definition")
}

func (s *Server) References(_ context.Context, _ *protocol.ReferenceParams) ([]protocol.Location, error) {
	return nil, notImplemented("References")
}

func (s *Server) DocumentHighlight(_ context.Context, _ *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	return nil, notImplemented("DocumentHighlight")
}

func (s *Server) DocumentSymbol(_ context.Context, _ *protocol.DocumentSymbolParams) ([]protocol.DocumentSymbol, error) {
	return nil, notImplemented("DocumentSymbol")
}

func (s *Server) CodeAction(_ context.Context, _ *protocol.CodeActionParams) ([]protocol.CodeAction, error) {
	return nil, notImplemented("CodeAction")
}

func (s *Server) Symbol(_ context.Context, _ *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	return nil, notImplemented("Symbol")
}

func (s *Server) CodeLens(_ context.Context, _ *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	return nil, notImplemented("CodeLens")
}

func (s *Server) ResolveCodeLens(_ context.Context, _ *protocol.CodeLens) (*protocol.CodeLens, error) {
	return nil, notImplemented("ResolveCodeLens")
}

func (s *Server) Formatting(_ context.Context, _ *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	return nil, notImplemented("Formatting")
}

func (s *Server) RangeFormatting(_ context.Context, _ *protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, notImplemented("RangeFormatting")
}

func (s *Server) OnTypeFormatting(_ context.Context, _ *protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, notImplemented("OnTypeFormatting")
}

func (s *Server) Rename(_ context.Context, _ *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	return nil, notImplemented("Rename")
}

func (s *Server) PrepareRename(_ context.Context, _ *protocol.PrepareRenameParams) (*protocol.Range, error) {
	return nil, notImplemented("PrepareRename")
}

func (s *Server) DocumentLink(_ context.Context, _ *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	return nil, notImplemented("DocumentLink")
}

func (s *Server) ResolveDocumentLink(_ context.Context, _ *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	return nil, notImplemented("ResolveDocumentLink")
}

func (s *Server) ExecuteCommand(_ context.Context, _ *protocol.ExecuteCommandParams) (interface{}, error) {
	return nil, notImplemented("ExecuteCommand")
}
