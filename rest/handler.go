// Copyright 2020 Tobias Guggenmos
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.  // You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/protocol"
	"github.com/prometheus-community/promql-langserver/langserver"
)

// CreateHandler creates an http.Handler for the PromQL langserver REST API.
//
// Expects the URL of a Prometheus server as the second argument.
func CreateHandler(ctx context.Context, prometheusURL string) (http.Handler, error) {
	langserver, err := langserver.CreateHeadlessServer(ctx, prometheusURL)
	if err != nil {
		return nil, err
	}

	return &langserverHandler{langserver: langserver}, nil
}

type langserverHandler struct {
	langserver     langserver.HeadlessServer
	requestCounter int64
}

func (h *langserverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var subHandler func(w http.ResponseWriter, r *http.Request, s langserver.HeadlessServer, requestID string)

	requestID := fmt.Sprint(atomic.AddInt64(&h.requestCounter, 1), ".promql")

	switch r.URL.Path {
	case "/diagnostics":
		subHandler = diagnosticsHandler
	case "/hover":
		subHandler = hoverHandler
	case "/completion":
		subHandler = completionHandler
	case "/signatureHelp":
		subHandler = signatureHelpHandler
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
		h.langserver.DidClose(r.Context(), &protocol.DidCloseTextDocumentParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: requestID,
			},
		},
		)
	}()

	if err := h.langserver.DidOpen(r.Context(), &protocol.DidOpenTextDocumentParams{
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

	subHandler(w, r, h.langserver, requestID)
}

func diagnosticsHandler(w http.ResponseWriter, r *http.Request, s langserver.HeadlessServer, requestID string) {
	hasLimit, limit, err := getLimitFromURL(r.URL)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	diagnostics, err := s.GetDiagnostics(requestID)
	if err != nil {
		http.Error(w, errors.Wrapf(err, "failed to get diagnostics").Error(), 500)
		return
	}

	items := diagnostics.Diagnostics

	if hasLimit && int64(len(items)) > limit {
		items = items[:limit]
	}

	returnJSON(w, items)
}

func hoverHandler(w http.ResponseWriter, r *http.Request, s langserver.HeadlessServer, requestID string) {
	position, err := getPositionFromURL(r.URL)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	hover, err := s.Hover(r.Context(), &protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: requestID,
			},
			Position: position,
		},
	})
	if err != nil {
		http.Error(w, errors.Wrapf(err, "failed to get hover info").Error(), 500)
		return
	}

	returnJSON(w, hover)
}

func completionHandler(w http.ResponseWriter, r *http.Request, s langserver.HeadlessServer, requestID string) {
	position, err := getPositionFromURL(r.URL)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	hasLimit, limit, err := getLimitFromURL(r.URL)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	completion, err := s.Completion(r.Context(), &protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: requestID,
			},
			Position: position,
		},
	})
	if err != nil {
		http.Error(w, errors.Wrapf(err, "failed to get completion info").Error(), 500)
		return
	}

	items := completion.Items

	if hasLimit && int64(len(items)) > limit {
		items = items[:limit]
	}

	returnJSON(w, items)
}

func signatureHelpHandler(w http.ResponseWriter, r *http.Request, s langserver.HeadlessServer, requestID string) {
	position, err := getPositionFromURL(r.URL)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	signature, err := s.SignatureHelp(r.Context(), &protocol.SignatureHelpParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: requestID,
			},
			Position: position,
		},
	})
	if err != nil {
		http.Error(w, errors.Wrapf(err, "failed to get hover info").Error(), 500)
		return
	}

	returnJSON(w, signature)
}

func returnJSON(w http.ResponseWriter, content interface{}) {
	encoder := json.NewEncoder(w)

	err := encoder.Encode(content)
	if err != nil {
		http.Error(w, errors.Wrapf(err, "failed to write response").Error(), 500)
	}
}

func getPositionFromURL(url *url.URL) (protocol.Position, error) {
	query := url.Query()
	lineStrs, ok := query["line"]

	if !ok || len(lineStrs) == 0 {
		return protocol.Position{}, errors.New("Param line is not specified")
	}

	line, err := strconv.ParseFloat(lineStrs[0], 64)
	if err != nil {
		return protocol.Position{}, errors.Wrap(err, "Failed to parse line number")
	}

	charStrs, ok := query["char"]

	if !ok || len(charStrs) == 0 {
		return protocol.Position{}, errors.New("Param char is not specified")
	}

	char, err := strconv.ParseFloat(charStrs[0], 64)
	if err != nil {
		return protocol.Position{}, errors.Wrap(err, "Failed to parse char number")
	}

	return protocol.Position{
		Line:      line,
		Character: char,
	}, nil
}

func getLimitFromURL(url *url.URL) (bool, int64, error) {
	query := url.Query()
	limitStrs, ok := query["limit"]

	if !ok || len(limitStrs) == 0 {
		return false, 0, nil
	}

	limit, err := strconv.ParseInt(limitStrs[0], 10, 64)
	if err != nil {
		return false, 0, errors.Wrap(err, "Failed to parse limit number")
	}

	if limit <= 0 {
		return false, 0, errors.New("Limit must be positive")
	}

	return true, limit, nil
}
