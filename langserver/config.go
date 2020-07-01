// Copyright 2019 The Prometheus Authors
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

	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/protocol"
)

// DidChangeConfiguration is required by the protocol.Server interface.
func (s *server) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error {
	langserverAddressConfigPath := []string{"promql", "url"}

	if params != nil {
		// nolint: errcheck
		s.client.LogMessage(
			s.lifetime,
			&protocol.LogMessageParams{
				Type:    protocol.Info,
				Message: fmt.Sprintf("Received notification change: %v\n", params),
			})

		setting := params.Settings

		for _, e := range langserverAddressConfigPath {
			m, ok := setting.(map[string]interface{})
			if !ok {
				break
			}

			setting, ok = m[e]
			if !ok {
				break
			}
		}

		if str, ok := setting.(string); ok {
			if err := s.connectPrometheus(str); err != nil {
				// nolint: errcheck
				s.client.LogMessage(ctx, &protocol.LogMessageParams{
					Type:    protocol.Info,
					Message: err.Error(),
				})
			}
		}
	}

	return nil
}

func (s *server) connectPrometheus(url string) error {
	if err := s.metadataService.ChangeDataSource(url); err != nil {
		// nolint: errcheck
		s.client.ShowMessage(s.lifetime, &protocol.ShowMessageParams{
			Type:    protocol.Error,
			Message: fmt.Sprintf("Failed to connect to Prometheus at %s:\n\n%s ", url, err.Error()),
		})
		return err
	}
	return nil
}
