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

	"github.com/slrtbtfs/go-tools-vendored/jsonrpc2"
	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
)

type Server struct {
	Conn   *jsonrpc2.Conn
	client protocol.Client

	state   serverState
	stateMu sync.Mutex

	cache documentCache
}

type serverState int

const (
	serverCreated = serverState(iota)
	serverInitializing
	serverInitialized // set once the server has received "initialized" request
	serverShutDown
)

func (s *Server) Run(_ context.Context) error {
	return s.Conn.Run(context.Background())
}

// Generates a Server from a jsonrpc2.Stream
func ServerFromStream(ctx context.Context, stream jsonrpc2.Stream) (context.Context, *Server) {
	s := &Server{}
	ctx, s.Conn, s.client = protocol.NewServer(ctx, stream, s)
	return ctx, s
}

// Generates a Server talking to stdio
func StdioServer(ctx context.Context) (context.Context, *Server) {
	stream := jsonrpc2.NewHeaderStream(os.Stdin, os.Stdout)
	stream = protocol.LoggingStream(stream, os.Stderr)
	return ServerFromStream(ctx, stream)
}
