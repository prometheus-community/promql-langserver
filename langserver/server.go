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

// This File includes code from the go/tools project which is governed by the following license:
// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package langserver

import (
	"context"
	"os"
	"sync"

	"github.com/slrtbtfs/promql-lsp/langserver/cache"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/jsonrpc2"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

// Server is a language server instance that can connect to exactly on client
type Server struct {
	Conn   *jsonrpc2.Conn
	client protocol.Client

	state   serverState
	stateMu sync.Mutex

	cache cache.DocumentCache

	config *Config
}

type serverState int

const (
	serverCreated = serverState(iota)
	serverInitializing
	serverInitialized // set once the server has received "initialized" request
	serverShutDown
)

// Run starts the language server instance
func (s *Server) Run(_ context.Context) error {
	return s.Conn.Run(context.Background())
}

// ServerFromStream generates a Server from a jsonrpc2.Stream
func ServerFromStream(ctx context.Context, stream jsonrpc2.Stream, config *Config) (context.Context, *Server) {
	s := &Server{}

	if config.Trace.Stderr {
		stream = protocol.LoggingStream(stream, os.Stderr)
	}

	ctx, s.Conn, s.client = protocol.NewServer(ctx, stream, s)
	s.config = config

	return ctx, s
}

// StdioServer generates a Server talking to stdio
func StdioServer(ctx context.Context, config *Config) (context.Context, *Server) {
	stream := jsonrpc2.NewHeaderStream(os.Stdin, os.Stdout)
	return ServerFromStream(ctx, stream, config)
}
