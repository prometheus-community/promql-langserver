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
	"io"
	"os"
	"time"

	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/jsonrpc2"
	"k8s.io/apimachinery/pkg/util/json"
)

type jsonLogStream struct {
	mainStream jsonrpc2.Stream
	logWriter  io.Writer
}

// jSONLogStream returns a stream that does log all communications in a format that
// can be streamed into the LSP inspector.
func jSONLogStream(str jsonrpc2.Stream, w io.Writer) jsonrpc2.Stream {
	ret := &jsonLogStream{str, w}

	return ret
}

func (s *jsonLogStream) Read(ctx context.Context) (jsonrpc2.Message, int64, error) {
	data, count, err := s.mainStream.Read(ctx)
	s.log(data, true)

	return data, count, err
}

func (s *jsonLogStream) Write(ctx context.Context, msg jsonrpc2.Message) (int64, error) {
	count, err := s.mainStream.Write(ctx, msg)
	s.log(msg, false)

	return count, err
}

func (s *jsonLogStream) Close() error {
	return s.mainStream.Close()
}

func getType(msg jsonrpc2.Message, incoming bool) (string, error) {
	var msgType string

	if incoming {
		msgType = "send-"
	} else {
		msgType = "receive-"
	}

	switch msg.(type) {
	case *jsonrpc2.Call:
		msgType += "request"
	case *jsonrpc2.Response:
		msgType += "response"
	case *jsonrpc2.Notification:
		msgType += "notification"
	}

	return msgType, nil
}

func (s *jsonLogStream) log(msg jsonrpc2.Message, incoming bool) {
	typ, err := getType(msg, incoming)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	msgText, err := json.Marshal(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	timestamp := time.Now().UnixNano() / 1000000
	tmformat := time.Now().Format("03:04:15.000 PM")
	// The LSP inspector expects the [LSP - <time>] part to be exactly 21 bytes
	fmt.Fprintf(s.logWriter, `[LSP-%s] {"type":"%s","message":%s,"timestamp":%d}%s`,
		tmformat, typ, msgText, timestamp, " \r\n")
}
