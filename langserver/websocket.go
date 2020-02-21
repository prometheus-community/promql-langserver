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

//nolint
import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	//"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/jsonrpc2"
)

// Implements the jsonrpc2.Stream interface
type wsConn struct {
	*websocket.Conn
}

func (c wsConn) Read(ctx context.Context) ([]byte, int64, error) {
	// Returning an error on an expired context is important here.
	// If that isn't done, the server won't stop after the connection
	// has been closed.
	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
	}

	typ, ret, err := c.ReadMessage()
	if err != nil {
		return nil, 0, err
	}

	if typ != websocket.TextMessage {
		return nil, 0, errors.New("wrong message type")
	}

	return ret, int64(len(ret)), nil
}

func (c wsConn) Write(ctx context.Context, msg []byte) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	err := c.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return 0, err
	}

	return int64(len(msg)), nil
}

func WebSocketHandler(addr string) (func(http.ResponseWriter, *http.Request), error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// No special handling required here, since the
			// error is already added to the ResponseWriter
			// by the upgraded call.
			return
		}

		ctx, cancel := context.WithCancel(context.TODO())

		ch := func(_ int, _ string) error {
			cancel()
			return nil
		}

		ws.SetCloseHandler(ch)

		var s Server

		_, s = ServerFromStream(ctx, wsConn{ws}, &Config{})

		if err := s.Run(); err != nil {
			// If the client disconnects, the above will fail
			return
		}
	}, nil
}
