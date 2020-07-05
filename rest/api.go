// Copyright 2020 The Prometheus Authors
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

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/protocol"
	"github.com/prometheus-community/promql-langserver/langserver"
	promClient "github.com/prometheus-community/promql-langserver/prometheus"
	"github.com/prometheus/common/route"
)

func returnJSON(w http.ResponseWriter, content interface{}) {
	encoder := json.NewEncoder(w)

	err := encoder.Encode(content)
	if err != nil {
		http.Error(w, errors.Wrapf(err, "failed to write response").Error(), 500)
	}
}

type lspData struct {
	// Expr is the promQL expression and is required for all endpoint available.
	Expr string `json:"expr"`
	// Limit is the number max of result returned to the client. It will be used for the autocompletion and the diagnostics.
	Limit *uint64 `json:"limit,omitempty"`
	// PositionLine is the number of the line for which the metadata is queried.
	PositionLine *float64 `json:"positionLine,omitempty"`
	// PositionChar for which the metadata is queried. Characters are counted as UTF16 Codepoints.
	PositionChar *float64 `json:"positionChar,omitempty"`
}

func (d *lspData) UnmarshalJSON(data []byte) error {
	var tmp lspData
	type plain lspData
	if err := json.Unmarshal(data, (*plain)(&tmp)); err != nil {
		return err
	}
	if len(tmp.Expr) == 0 {
		return fmt.Errorf("promQL expression is not specified")
	}
	*d = tmp
	return nil
}

func (d *lspData) returnPosition() (protocol.Position, error) {
	if d.PositionLine == nil {
		return protocol.Position{}, errors.New("positionLine is not specified")
	}
	if d.PositionChar == nil {
		return protocol.Position{}, errors.New("positionChar is not specified")
	}
	return protocol.Position{
		Line:      *d.PositionLine,
		Character: *d.PositionChar,
	}, nil
}

type API struct {
	langServer    langserver.HeadlessServer
	enableMetrics bool
}

func NewLangServerAPI(ctx context.Context, metadataService promClient.MetadataService, logger log.Logger, enableMetrics bool) (*API, error) {
	lgs, err := langserver.CreateHeadlessServer(ctx, metadataService, logger)
	if err != nil {
		return nil, err
	}
	return &API{
		langServer:    lgs,
		enableMetrics: enableMetrics,
	}, nil
}

// Register the API's endpoints in the given router.
func (a *API) Register(r *route.Router, prefix string) {
	r.Post(prefix+"/diagnostics", a.handle(a.diagnostics))
	r.Post(prefix+"/completion", a.handle(a.completion))
	r.Post(prefix+"/hover", a.handle(a.hover))
	r.Post(prefix+"/signatureHelp", a.handle(a.signature))
}

func (a *API) handle(h http.HandlerFunc) http.HandlerFunc {
	endpoint := manageDocumentMiddleware(a.langServer)(h)
	return endpoint
}

func (a *API) diagnostics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID, requestData, err := getRequestDataAndID(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	diagnostics, err := a.langServer.GetDiagnostics(requestID)
	if err != nil {
		http.Error(w, errors.Wrapf(err, "failed to get diagnostics").Error(), 500)
		return
	}

	items := diagnostics.Diagnostics
	limit := requestData.Limit
	if limit != nil && uint64(len(items)) > *limit {
		items = items[:*limit]
	}

	returnJSON(w, items)
}

func (a *API) hover(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID, requestData, err := getRequestDataAndID(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	position, err := requestData.returnPosition()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	hover, err := a.langServer.Hover(r.Context(), &protocol.HoverParams{
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

func (a *API) completion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID, requestData, err := getRequestDataAndID(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	position, err := requestData.returnPosition()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	completion, err := a.langServer.Completion(r.Context(), &protocol.CompletionParams{
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
	limit := requestData.Limit
	if limit != nil && uint64(len(items)) > *limit {
		items = items[:*limit]
	}

	returnJSON(w, items)
}

func (a *API) signature(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID, requestData, err := getRequestDataAndID(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	position, err := requestData.returnPosition()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	signature, err := a.langServer.SignatureHelp(r.Context(), &protocol.SignatureHelpParams{
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
