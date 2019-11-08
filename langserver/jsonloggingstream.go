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

	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/jsonrpc2"
)

type jsonLogStream struct {
	stream  jsonrpc2.Stream
	log     io.Writer
	waiting chan logItem
}

type logItem struct {
	msg      []byte
	incoming bool
}

// JSONLogStream returns a stream that does log all communications in a format that
// can be streamed into the lsp inspector
func JSONLogStream(str jsonrpc2.Stream, w io.Writer) jsonrpc2.Stream {
	ret := &jsonLogStream{str, w, make(chan logItem)}
	go ret.startLogging()

	return ret
}

func (s *jsonLogStream) Read(ctx context.Context) ([]byte, int64, error) {
	data, count, err := s.stream.Read(ctx)
	s.waiting <- logItem{data, true}

	return data, count, err
}

func (s *jsonLogStream) Write(ctx context.Context, data []byte) (int64, error) {
	count, err := s.stream.Write(ctx, data)
	s.waiting <- logItem{data, false}

	return count, err
}

func getType(msg []byte, incoming bool) (string, error) {
	var v protocol.Combined

	var msgType string

	err := json.Unmarshal(msg, &v)
	if err != nil {
		return "", err
	}

	if incoming {
		msgType = "send-"
	} else {
		msgType = "receive-"
	}

	switch {
	case v.ID != nil && v.Method != "" && (v.Params != nil || v.Method == "shutdown"):
		msgType += "request"
	case v.ID != nil && v.Method == "" && v.Params == nil:
		msgType += "response"
	default:
		// This might be not always accurate
		msgType += "notification"
	}

	return msgType, nil
}

func (s *jsonLogStream) startLogging() {
	for item := range s.waiting {
		typ, err := getType(item.msg, item.incoming)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		timestamp := time.Now().UnixNano() / 1000000
		tmformat := time.Now().Format("03:04:15.000 PM")
		fmt.Fprintf(s.log, `[LSP-%s] {"isLSPMessage":true,"type":"%s","message":%s,"timestamp":%d}%s`,
			tmformat, typ, item.msg, timestamp, " \r\n")
	}
}
