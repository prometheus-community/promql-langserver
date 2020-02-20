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

package stateless

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/prometheus-community/promql-langserver/langserver"
	"github.com/prometheus-community/promql-langserver/vendored/go-tools/lsp/protocol"
)

// Create an API handler for the stateless langserver API
// Expects the URL of a Prometheus server as the argument.
// Will fail if the Prometheus server is not reachable
func CreateAPIHandler(ctx context.Context, prometheusURL string) (http.Handler, error) {
	langserver, err := langserver.CreateHeadlessServer(ctx, prometheusURL)
	if err != nil {
		return nil, err
	}

	return &langserverHandler{ctx: ctx, langserver: langserver}, nil
}

type langserverHandler struct {
	langserver     langserver.HeadlessServer
	requestCounter int64
	ctx            context.Context
}

func (h *langserverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, r.URL.Path)

	var subHandler func(http.ResponseWriter, *http.Request)

	requestID := fmt.Sprint(atomic.AddInt64(&h.requestCounter, 1), ".promql")

	switch r.URL.Path {
	case "/diagnostics":
		subHandler = diagnosticsHandler(h.langserver, requestID)
	default:
		http.NotFound(w, r)
		return
	}

	exprs, ok := r.URL.Query()["expr"]

	if !ok || len(exprs) == 0 {
		http.Error(w, "Param expr is not specified", 400)
		return
	}

	defer func() {
		h.langserver.DidClose(h.ctx, &protocol.DidCloseTextDocumentParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: requestID,
			},
		},
		)
	}()

	if err := h.langserver.DidOpen(h.ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        requestID,
			LanguageID: "promql",
			Version:    0,
			Text:       exprs[0],
		},
	}); err != nil {
		http.Error(w, errors.Wrapf(err, "failed to open document").Error(), 500)
		return
	}

	subHandler(w, r)

}

func diagnosticsHandler(s langserver.HeadlessServer, uri string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		diagnostics, err := s.GetDiagnostics(uri)
		if err != nil {
			http.Error(w, errors.Wrapf(err, "failed to get diagnostics").Error(), 500)
			return
		}

		returnJSON(w, diagnostics.Diagnostics)

	}
}

func returnJSON(w http.ResponseWriter, content interface{}) {
	encoder := json.NewEncoder(w)

	err := encoder.Encode(content)
	if err != nil {
		http.Error(w, errors.Wrapf(err, "failed to write response").Error(), 500)
	}
}
