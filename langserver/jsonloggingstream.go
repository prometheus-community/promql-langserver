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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"go.lsp.dev/jsonrpc2"
)

type jsonLogStream struct {
	mainStream jsonrpc2.Stream
	logWriter  io.Writer
}

// jSONLogStream returns a stream that does log all communications in a format that
// can be streamed into the LSP inspector.
func jSONLogStream(str jsonrpc2.Stream, w io.Writer) jsonrpc2.Stream {
	return &jsonLogStream{str, w}
}

func (s *jsonLogStream) Read(ctx context.Context) (jsonrpc2.Message, int64, error) {
	msg, count, err := s.mainStream.Read(ctx)
	if err == nil {
		s.log(msg, true)
	}

	return msg, count, err
}

func (s *jsonLogStream) Write(ctx context.Context, msg jsonrpc2.Message) (int64, error) {
	s.log(msg, false)

	return s.mainStream.Write(ctx, msg)
}

func (s *jsonLogStream) Close() error {
	return s.mainStream.Close()
}

// messageType classifies a JSON-RPC message for the LSP inspector log.
func messageType(msg jsonrpc2.Message, incoming bool) string {
	direction := "receive-"
	if incoming {
		direction = "send-"
	}

	switch msg.(type) {
	case *jsonrpc2.Call:
		return direction + "request"
	case *jsonrpc2.Response:
		return direction + "response"
	default:
		return direction + "notification"
	}
}

func (s *jsonLogStream) log(msg jsonrpc2.Message, incoming bool) {
	if msg == nil {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	timestamp := time.Now().UnixNano() / 1000000
	tmformat := time.Now().Format("03:04:15.000 PM")
	// The LSP inspector expects the [LSP - <time>] part to be exactly 21 bytes
	fmt.Fprintf(s.logWriter, `[LSP-%s] {"type":"%s","message":%s,"timestamp":%d}%s`,
		tmformat, messageType(msg, incoming), data, timestamp, " \r\n")
}
