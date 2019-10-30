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
	"fmt"
	"os"

	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/jsonrpc2"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

func notImplemented(method string) *jsonrpc2.Error {
	err := jsonrpc2.NewErrorf(jsonrpc2.CodeMethodNotFound, "method %q no yet implemented", method)
	fmt.Fprint(os.Stderr, err.Error())

	return err
}

// DidChangeWorkspaceFolders is required by the protocol.Server interface
func (s *Server) DidChangeWorkspaceFolders(_ context.Context, _ *protocol.DidChangeWorkspaceFoldersParams) error {
	return notImplemented("DidChangeWorkspaceFolders")
}

// DidChangeConfiguration is required by the protocol.Server interface
func (s *Server) DidChangeConfiguration(_ context.Context, _ *protocol.DidChangeConfigurationParams) error {
	//return notImplemented("DidChangeConfiguration")
	// For ycmd
	return nil
}

// DidSave is required by the protocol.Server interface
func (s *Server) DidSave(_ context.Context, _ *protocol.DidSaveTextDocumentParams) error {
	return notImplemented("DidSave")
}

// WillSave is required by the protocol.Server interface
func (s *Server) WillSave(_ context.Context, _ *protocol.WillSaveTextDocumentParams) error {
	return notImplemented("WillSave")
}

// DidChangeWatchedFiles is required by the protocol.Server interface
func (s *Server) DidChangeWatchedFiles(_ context.Context, _ *protocol.DidChangeWatchedFilesParams) error {
	return notImplemented("DidChangeWatchedFiles")
}

// Progress is required by the protocol.Server interface
func (s *Server) Progress(_ context.Context, _ *protocol.ProgressParams) error {
	return notImplemented("Progress")
}

// SelectionRange is required by the protocol.Server interface
func (s *Server) SelectionRange(_ context.Context, _ *protocol.SelectionRangeParams) ([]protocol.SelectionRange, error) { // nolint:lll
	return nil, notImplemented("SelectionRange")
}

// SetTraceNotification is required by the protocol.Server interface
func (s *Server) SetTraceNotification(_ context.Context, _ *protocol.SetTraceParams) error {
	return notImplemented("SetTraceNotification")
}

// LogTraceNotification is required by the protocol.Server interface
func (s *Server) LogTraceNotification(_ context.Context, _ *protocol.LogTraceParams) error {
	return notImplemented("LogTraceNotification")
}

// Implementation is required by the protocol.Server interface
func (s *Server) Implementation(_ context.Context, _ *protocol.ImplementationParams) ([]protocol.Location, error) {
	return nil, notImplemented("Implementation")
}

// TypeDefinition is required by the protocol.Server interface
func (s *Server) TypeDefinition(_ context.Context, _ *protocol.TypeDefinitionParams) ([]protocol.Location, error) {
	return nil, notImplemented("TypeDefinition")
}

// DocumentColor is required by the protocol.Server interface
func (s *Server) DocumentColor(_ context.Context, _ *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) { //nolint:lll
	return nil, notImplemented("DocumentColor")
}

// ColorPresentation is required by the protocol.Server interface
func (s *Server) ColorPresentation(_ context.Context, _ *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) { //nolint:lll
	return nil, notImplemented("ColorPresentation")
}

// FoldingRange is required by the protocol.Server interface
func (s *Server) FoldingRange(_ context.Context, _ *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	return nil, notImplemented("FoldingRange")
}

// Declaration is required by the protocol.Server interface
func (s *Server) Declaration(_ context.Context, _ *protocol.DeclarationParams) ([]protocol.DeclarationLink, error) {
	return nil, notImplemented("Declaration")
}

// WillSaveWaitUntil is required by the protocol.Server interface
func (s *Server) WillSaveWaitUntil(_ context.Context, _ *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) { //nolint:lll
	return nil, notImplemented("WillSaveWaitUntil")
}

// Completion is required by the protocol.Server interface
// nolint: wsl
func (s *Server) Completion(_ context.Context, _ *protocol.CompletionParams) (*protocol.CompletionList, error) {
	// For ycmd
	/*
		return &protocol.CompletionList{
			IsIncomplete: true,
			Items:        nil,
		}, nil
	*/
	return nil, notImplemented("Completion")
}

// Resolve is required by the protocol.Server interface
func (s *Server) Resolve(_ context.Context, _ *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return nil, notImplemented("Resolve")
}

// SignatureHelp is required by the protocol.Server interface
func (s *Server) SignatureHelp(_ context.Context, _ *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	return nil, notImplemented("SignatureHelp")
}

// Definition is required by the protocol.Server interface
func (s *Server) Definition(_ context.Context, _ *protocol.DefinitionParams) ([]protocol.Location, error) {
	return nil, notImplemented("Definition")
}

// References is required by the protocol.Server interface
func (s *Server) References(_ context.Context, _ *protocol.ReferenceParams) ([]protocol.Location, error) {
	return nil, notImplemented("References")
}

// DocumentHighlight is required by the protocol.Server interface
func (s *Server) DocumentHighlight(_ context.Context, _ *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) { //nolint:lll
	return nil, notImplemented("DocumentHighlight")
}

// DocumentSymbol is required by the protocol.Server interface
func (s *Server) DocumentSymbol(_ context.Context, _ *protocol.DocumentSymbolParams) ([]protocol.DocumentSymbol, error) { //nolint:lll
	return nil, notImplemented("DocumentSymbol")
}

// CodeAction is required by the protocol.Server interface
func (s *Server) CodeAction(_ context.Context, _ *protocol.CodeActionParams) ([]protocol.CodeAction, error) {
	return nil, notImplemented("CodeAction")
}

// Symbol is required by the protocol.Server interface
func (s *Server) Symbol(_ context.Context, _ *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	return nil, notImplemented("Symbol")
}

// CodeLens is required by the protocol.Server interface
func (s *Server) CodeLens(_ context.Context, _ *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	return nil, notImplemented("CodeLens")
}

// ResolveCodeLens is required by the protocol.Server interface
func (s *Server) ResolveCodeLens(_ context.Context, _ *protocol.CodeLens) (*protocol.CodeLens, error) {
	return nil, notImplemented("ResolveCodeLens")
}

// Formatting is required by the protocol.Server interface
func (s *Server) Formatting(_ context.Context, _ *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	return nil, notImplemented("Formatting")
}

// RangeFormatting is required by the protocol.Server interface
func (s *Server) RangeFormatting(_ context.Context, _ *protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit, error) { //nolint:lll
	return nil, notImplemented("RangeFormatting")
}

// OnTypeFormatting is required by the protocol.Server interface
func (s *Server) OnTypeFormatting(_ context.Context, _ *protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit, error) { // nolint:lll
	return nil, notImplemented("OnTypeFormatting")
}

// Rename is required by the protocol.Server interface
func (s *Server) Rename(_ context.Context, _ *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	return nil, notImplemented("Rename")
}

// PrepareRename is required by the protocol.Server interface
func (s *Server) PrepareRename(_ context.Context, _ *protocol.PrepareRenameParams) (*protocol.Range, error) {
	return nil, notImplemented("PrepareRename")
}

// DocumentLink is required by the protocol.Server interface
func (s *Server) DocumentLink(_ context.Context, _ *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	return nil, notImplemented("DocumentLink")
}

// ResolveDocumentLink is required by the protocol.Server interface
func (s *Server) ResolveDocumentLink(_ context.Context, _ *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	return nil, notImplemented("ResolveDocumentLink")
}

// ExecuteCommand is required by the protocol.Server interface
func (s *Server) ExecuteCommand(_ context.Context, _ *protocol.ExecuteCommandParams) (interface{}, error) {
	return nil, notImplemented("ExecuteCommand")
}
