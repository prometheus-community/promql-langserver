// Copyright 2020 Tobias Guggenmos
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

	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/protocol"
)

type headlessClient struct{}

func (headlessClient) ShowMessage(_ context.Context, params *protocol.ShowMessageParams) error {
	fmt.Fprintln(os.Stderr, params.Message)
	return nil
}

func (headlessClient) LogMessage(_ context.Context, params *protocol.LogMessageParams) error {
	fmt.Fprintln(os.Stderr, params.Message)
	return nil
}

func (headlessClient) Event(_ context.Context, _ *interface{}) error {
	// ignore
	return nil
}

func (headlessClient) PublishDiagnostics(_ context.Context, _ *protocol.PublishDiagnosticsParams) error {
	// ignore
	return nil
}

func (headlessClient) WorkspaceFolders(_ context.Context) ([]protocol.WorkspaceFolder, error) {
	// ignore
	return nil, nil
}

func (headlessClient) Configuration(_ context.Context, _ *protocol.ParamConfiguration) ([]interface{}, error) {
	// ignore
	return nil, nil
}

func (headlessClient) RegisterCapability(_ context.Context, _ *protocol.RegistrationParams) error {
	// ignore
	return nil
}

func (headlessClient) UnregisterCapability(_ context.Context, _ *protocol.UnregistrationParams) error {
	// ignore
	return nil
}

func (headlessClient) ShowMessageRequest(_ context.Context, _ *protocol.ShowMessageRequestParams) (*protocol.MessageActionItem, error) {
	// ignore
	return nil, nil
}

func (headlessClient) ApplyEdit(_ context.Context, _ *protocol.ApplyWorkspaceEditParams) (*protocol.ApplyWorkspaceEditResponse, error) {
	// ignore
	return nil, nil
}
