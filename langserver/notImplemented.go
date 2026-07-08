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

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

func notImplemented(method string) *jsonrpc2.Error {
	return jsonrpc2.Errorf(jsonrpc2.MethodNotFound, "method %q not yet implemented", method)
}

// WorkDoneProgressCancel is required by the protocol.Server interface.
func (s *server) WorkDoneProgressCancel(_ context.Context, _ *protocol.WorkDoneProgressCancelParams) error {
	return notImplemented("WorkDoneProgressCancel")
}

// LogTrace is required by the protocol.Server interface.
func (s *server) LogTrace(_ context.Context, _ *protocol.LogTraceParams) error {
	return notImplemented("LogTrace")
}

// SetTrace is required by the protocol.Server interface.
func (s *server) SetTrace(_ context.Context, _ *protocol.SetTraceParams) error {
	return notImplemented("SetTrace")
}

// CodeAction is required by the protocol.Server interface.
func (s *server) CodeAction(_ context.Context, _ *protocol.CodeActionParams) ([]protocol.CodeAction, error) {
	return nil, notImplemented("CodeAction")
}

// CodeLens is required by the protocol.Server interface.
func (s *server) CodeLens(_ context.Context, _ *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	// The language client cannot be told to stop asking for Code Lenses, so to
	// prevent the editor from showing error messages this is implemented by
	// returning an empty result.
	return nil, nil
}

// CodeLensResolve is required by the protocol.Server interface.
func (s *server) CodeLensResolve(_ context.Context, _ *protocol.CodeLens) (*protocol.CodeLens, error) {
	return nil, notImplemented("CodeLensResolve")
}

// CodeLensRefresh is required by the protocol.Server interface.
func (s *server) CodeLensRefresh(_ context.Context) error {
	return notImplemented("CodeLensRefresh")
}

// ColorPresentation is required by the protocol.Server interface.
func (s *server) ColorPresentation(_ context.Context, _ *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	return nil, notImplemented("ColorPresentation")
}

// CompletionResolve is required by the protocol.Server interface.
func (s *server) CompletionResolve(_ context.Context, _ *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return nil, notImplemented("CompletionResolve")
}

// Declaration is required by the protocol.Server interface.
func (s *server) Declaration(_ context.Context, _ *protocol.DeclarationParams) ([]protocol.Location, error) {
	return nil, notImplemented("Declaration")
}

// DidChangeWatchedFiles is required by the protocol.Server interface.
func (s *server) DidChangeWatchedFiles(_ context.Context, _ *protocol.DidChangeWatchedFilesParams) error {
	return notImplemented("DidChangeWatchedFiles")
}

// DidChangeWorkspaceFolders is required by the protocol.Server interface.
func (s *server) DidChangeWorkspaceFolders(_ context.Context, _ *protocol.DidChangeWorkspaceFoldersParams) error {
	return notImplemented("DidChangeWorkspaceFolders")
}

// DidSave is required by the protocol.Server interface.
func (s *server) DidSave(_ context.Context, _ *protocol.DidSaveTextDocumentParams) error {
	return notImplemented("DidSave")
}

// DocumentColor is required by the protocol.Server interface.
func (s *server) DocumentColor(_ context.Context, _ *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	return nil, notImplemented("DocumentColor")
}

// DocumentHighlight is required by the protocol.Server interface.
func (s *server) DocumentHighlight(_ context.Context, _ *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	return nil, notImplemented("DocumentHighlight")
}

// DocumentLink is required by the protocol.Server interface.
func (s *server) DocumentLink(_ context.Context, _ *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	// The language client cannot be told to stop asking for Document Links, so
	// to prevent the editor from showing error messages this is implemented by
	// returning an empty result.
	return nil, nil
}

// DocumentLinkResolve is required by the protocol.Server interface.
func (s *server) DocumentLinkResolve(_ context.Context, _ *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	return nil, notImplemented("DocumentLinkResolve")
}

// DocumentSymbol is required by the protocol.Server interface.
func (s *server) DocumentSymbol(_ context.Context, _ *protocol.DocumentSymbolParams) ([]interface{}, error) {
	return nil, notImplemented("DocumentSymbol")
}

// ExecuteCommand is required by the protocol.Server interface.
func (s *server) ExecuteCommand(_ context.Context, _ *protocol.ExecuteCommandParams) (interface{}, error) {
	return nil, notImplemented("ExecuteCommand")
}

// FoldingRanges is required by the protocol.Server interface.
func (s *server) FoldingRanges(_ context.Context, _ *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	return nil, notImplemented("FoldingRanges")
}

// Formatting is required by the protocol.Server interface.
func (s *server) Formatting(_ context.Context, _ *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	return nil, notImplemented("Formatting")
}

// Implementation is required by the protocol.Server interface.
func (s *server) Implementation(_ context.Context, _ *protocol.ImplementationParams) ([]protocol.Location, error) {
	return nil, notImplemented("Implementation")
}

// OnTypeFormatting is required by the protocol.Server interface.
func (s *server) OnTypeFormatting(_ context.Context, _ *protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, notImplemented("OnTypeFormatting")
}

// PrepareRename is required by the protocol.Server interface.
func (s *server) PrepareRename(_ context.Context, _ *protocol.PrepareRenameParams) (*protocol.Range, error) {
	return nil, notImplemented("PrepareRename")
}

// RangeFormatting is required by the protocol.Server interface.
func (s *server) RangeFormatting(_ context.Context, _ *protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, notImplemented("RangeFormatting")
}

// References is required by the protocol.Server interface.
func (s *server) References(_ context.Context, _ *protocol.ReferenceParams) ([]protocol.Location, error) {
	return nil, notImplemented("References")
}

// Rename is required by the protocol.Server interface.
func (s *server) Rename(_ context.Context, _ *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	return nil, notImplemented("Rename")
}

// Symbols is required by the protocol.Server interface.
func (s *server) Symbols(_ context.Context, _ *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	return nil, notImplemented("Symbols")
}

// TypeDefinition is required by the protocol.Server interface.
func (s *server) TypeDefinition(_ context.Context, _ *protocol.TypeDefinitionParams) ([]protocol.Location, error) {
	return nil, notImplemented("TypeDefinition")
}

// WillSave is required by the protocol.Server interface.
func (s *server) WillSave(_ context.Context, _ *protocol.WillSaveTextDocumentParams) error {
	return notImplemented("WillSave")
}

// WillSaveWaitUntil is required by the protocol.Server interface.
func (s *server) WillSaveWaitUntil(_ context.Context, _ *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	return nil, notImplemented("WillSaveWaitUntil")
}

// ShowDocument is required by the protocol.Server interface.
func (s *server) ShowDocument(_ context.Context, _ *protocol.ShowDocumentParams) (*protocol.ShowDocumentResult, error) {
	return nil, notImplemented("ShowDocument")
}

// WillCreateFiles is required by the protocol.Server interface.
func (s *server) WillCreateFiles(_ context.Context, _ *protocol.CreateFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, notImplemented("WillCreateFiles")
}

// DidCreateFiles is required by the protocol.Server interface.
func (s *server) DidCreateFiles(_ context.Context, _ *protocol.CreateFilesParams) error {
	return notImplemented("DidCreateFiles")
}

// WillRenameFiles is required by the protocol.Server interface.
func (s *server) WillRenameFiles(_ context.Context, _ *protocol.RenameFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, notImplemented("WillRenameFiles")
}

// DidRenameFiles is required by the protocol.Server interface.
func (s *server) DidRenameFiles(_ context.Context, _ *protocol.RenameFilesParams) error {
	return notImplemented("DidRenameFiles")
}

// WillDeleteFiles is required by the protocol.Server interface.
func (s *server) WillDeleteFiles(_ context.Context, _ *protocol.DeleteFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, notImplemented("WillDeleteFiles")
}

// DidDeleteFiles is required by the protocol.Server interface.
func (s *server) DidDeleteFiles(_ context.Context, _ *protocol.DeleteFilesParams) error {
	return notImplemented("DidDeleteFiles")
}

// PrepareCallHierarchy is required by the protocol.Server interface.
func (s *server) PrepareCallHierarchy(_ context.Context, _ *protocol.CallHierarchyPrepareParams) ([]protocol.CallHierarchyItem, error) {
	return nil, notImplemented("PrepareCallHierarchy")
}

// IncomingCalls is required by the protocol.Server interface.
func (s *server) IncomingCalls(_ context.Context, _ *protocol.CallHierarchyIncomingCallsParams) ([]protocol.CallHierarchyIncomingCall, error) {
	return nil, notImplemented("IncomingCalls")
}

// OutgoingCalls is required by the protocol.Server interface.
func (s *server) OutgoingCalls(_ context.Context, _ *protocol.CallHierarchyOutgoingCallsParams) ([]protocol.CallHierarchyOutgoingCall, error) {
	return nil, notImplemented("OutgoingCalls")
}

// SemanticTokensFull is required by the protocol.Server interface.
func (s *server) SemanticTokensFull(_ context.Context, _ *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	return nil, notImplemented("SemanticTokensFull")
}

// SemanticTokensFullDelta is required by the protocol.Server interface.
func (s *server) SemanticTokensFullDelta(_ context.Context, _ *protocol.SemanticTokensDeltaParams) (interface{}, error) {
	return nil, notImplemented("SemanticTokensFullDelta")
}

// SemanticTokensRange is required by the protocol.Server interface.
func (s *server) SemanticTokensRange(_ context.Context, _ *protocol.SemanticTokensRangeParams) (*protocol.SemanticTokens, error) {
	return nil, notImplemented("SemanticTokensRange")
}

// SemanticTokensRefresh is required by the protocol.Server interface.
func (s *server) SemanticTokensRefresh(_ context.Context) error {
	return notImplemented("SemanticTokensRefresh")
}

// LinkedEditingRange is required by the protocol.Server interface.
func (s *server) LinkedEditingRange(_ context.Context, _ *protocol.LinkedEditingRangeParams) (*protocol.LinkedEditingRanges, error) {
	return nil, notImplemented("LinkedEditingRange")
}

// Moniker is required by the protocol.Server interface.
func (s *server) Moniker(_ context.Context, _ *protocol.MonikerParams) ([]protocol.Moniker, error) {
	return nil, notImplemented("Moniker")
}

// Request is required by the protocol.Server interface.
func (s *server) Request(_ context.Context, _ string, _ interface{}) (interface{}, error) {
	return nil, notImplemented("Request")
}
