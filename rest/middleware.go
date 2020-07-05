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
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/protocol"
	"github.com/prometheus-community/promql-langserver/langserver"
)

type key int

const (
	contextKeyRequestID key = iota
	contextKeyRequestData
)

// injectRequestID will generate a new context with the requestID injected.
func injectRequestID(ctx context.Context, requestID protocol.DocumentURI) context.Context {
	return context.WithValue(ctx, contextKeyRequestID, requestID)
}

func getRequestID(ctx context.Context) (protocol.DocumentURI, error) {
	reqID := ctx.Value(contextKeyRequestID)
	if ret, ok := reqID.(protocol.DocumentURI); ok {
		return ret, nil
	}
	return "", fmt.Errorf("unable to retrieve the requestID")
}

// injectRequestData will create a new context and add to this new one the data passed as parameter.
func injectRequestData(ctx context.Context, data *lspData) context.Context {
	return context.WithValue(ctx, contextKeyRequestData, data)
}

func getRequestData(ctx context.Context) (*lspData, error) {
	reqData := ctx.Value(contextKeyRequestData)
	if ret, ok := reqData.(*lspData); ok {
		return ret, nil
	}
	return nil, fmt.Errorf("unable to retrieve the request data")
}

func getRequestDataAndID(ctx context.Context) (protocol.DocumentURI, *lspData, error) {
	id, err := getRequestID(ctx)
	if err != nil {
		return "", nil, err
	}
	data, err := getRequestData(ctx)
	if err != nil {
		return "", nil, err
	}
	return id, data, nil
}

type middlewareFunc func(http.HandlerFunc) http.HandlerFunc

// manageDocumentMiddleware is an HTTP middleware that will:
//   * generate and inject a requestID
//   * unmarshal the body and inject it in the request context
//   * open the document used to analyze
//   * assure that the document used to analyze will be closed properly at the end of the HTTP request.
func manageDocumentMiddleware(langServer langserver.HeadlessServer) middlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// start to generate an unique ID for the given request
			id, err := uuid.NewRandom()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			requestID := protocol.DocumentURI(id.String())
			// then inject it in the http request context
			r = r.WithContext(injectRequestID(r.Context(), requestID))

			// then unmarshall the body to the proper struct to be able to retrieve the promQL expr
			data := &lspData{}
			if err := json.NewDecoder(r.Body).Decode(data); err != nil {
				if err == io.EOF {
					// this case is used just in order to have a proper error message instead of just "EOF"
					http.Error(w, fmt.Sprint("body not present"), 400)
					return
				}
				http.Error(w, err.Error(), 400)
				return
			}
			// inject the data unmarshalled to avoid to have to decode it later
			r = r.WithContext(injectRequestData(r.Context(), data))

			// open the document to prepare the next operation on it
			if err := langServer.DidOpen(r.Context(), &protocol.DidOpenTextDocumentParams{
				TextDocument: protocol.TextDocumentItem{
					URI:        requestID,
					LanguageID: "promql",
					Version:    0,
					Text:       data.Expr,
				},
			}); err != nil {
				http.Error(w, errors.Wrapf(err, "failed to open document").Error(), 500)
				return
			}

			defer func() {
				// assure to close the document in order to stay stateless
				langServer.DidClose(r.Context(), &protocol.DidCloseTextDocumentParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: requestID,
					},
				},
				)
			}()
			next(w, r)
		}
	}
}
