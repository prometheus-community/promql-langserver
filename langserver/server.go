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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/blang/semver"
	"github.com/pkg/errors"

	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/jsonrpc2"
	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/protocol"
	"github.com/prometheus-community/promql-langserver/langserver/cache"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// Server wraps language server instance that can connect to exactly one client.
type Server struct {
	server *server
}

// HeadlessServer is a modified Server interface that is used by the REST API.
type HeadlessServer interface {
	protocol.Server
	GetDiagnostics(uri protocol.DocumentURI) (*protocol.PublishDiagnosticsParams, error)
}

// server is a language server instance that can connect to exactly one client
type server struct {
	Conn   *jsonrpc2.Conn
	client protocol.Client

	state   serverState
	stateMu sync.Mutex

	cache cache.DocumentCache

	config *Config

	prometheus        api.Client
	PrometheusURL     string
	PrometheusVersion string
	prometheusMu      sync.Mutex

	lifetime context.Context
	exit     func()
	headless bool
}

type serverState int

type testResponse struct {
	Status    string        `json:"status"`
	Data      buildInfoData `json:"data,omitempty"`
	ErrorType string        `json:"errorType,omitempty"`
	Error     string        `json:"error,omitempty"`
	Warnings  []string      `json:"warnings,omitempty"`
}

// buildInfoData contains build information about Prometheus.
type buildInfoData struct {
	Version   string `json:"version"`
	Revision  string `json:"revision"`
	Branch    string `json:"branch"`
	BuildUser string `json:"buildUser"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
}

const (
	serverCreated = serverState(iota)
	serverInitializing
	serverInitialized // set once the server has received "initialized" request
	serverShutDown
)

// Run starts the language server instance.
func (s Server) Run() error {
	return s.server.Conn.Run(s.server.lifetime)
}

// CreateHeadlessServer creates a locked down server instance for the REST API.
//
// "locked down" in this case means, that the instance cannot send or receive any JSONRPC communication. Logging messages that the instance tries to send over JSONRPC are redirected to stderr.
func CreateHeadlessServer(ctx context.Context, prometheusURL string) (HeadlessServer, error) {
	s := &server{
		client:   headlessClient{},
		headless: true,
		config:   &Config{PrometheusURL: prometheusURL},
	}

	s.lifetime, s.exit = context.WithCancel(ctx)

	if _, err := s.Initialize(ctx, &protocol.ParamInitialize{}); err != nil {
		return nil, err
	}

	if err := s.Initialized(ctx, &protocol.InitializedParams{}); err != nil {
		return nil, err
	}

	return s, nil
}

// ServerFromStream generates a Server from a jsonrpc2.Stream.
func ServerFromStream(ctx context.Context, stream jsonrpc2.Stream, config *Config) (context.Context, Server) {
	s := &server{}

	switch config.RPCTrace {
	case "text":
		stream = protocol.LoggingStream(stream, os.Stderr)
	case "json":
		stream = jSONLogStream(stream, os.Stderr)
	}

	s.Conn = jsonrpc2.NewConn(stream)
	s.client = protocol.ClientDispatcher(s.Conn)

	s.Conn.AddHandler(protocol.ServerHandler(s))

	s.config = config

	s.lifetime, s.exit = context.WithCancel(ctx)

	return ctx, Server{s}
}

// nolint:funlen
func (s *server) connectPrometheus(url string) error {
	s.prometheusMu.Lock()
	defer s.prometheusMu.Unlock()

	s.PrometheusURL = ""
	s.prometheus = nil

	if strings.TrimSpace(url) == "" {
		return nil
	}

	var err error

	s.prometheus, err = api.NewClient(api.Config{Address: url})
	err = errors.Wrapf(err, "Failed to connect to prometheus: %s\n", url)

	if err == nil {
		// nolint: errcheck
		s.client.LogMessage(s.lifetime, &protocol.LogMessageParams{
			Type:    protocol.Info,
			Message: fmt.Sprint("Prometheus: ", url),
		})
	}

	if !s.headless {
		testurl := fmt.Sprint(url, "/api/v1/status/buildinfo")

		resp, err := http.Get(testurl) // nolint: gosec

		if err != nil {
			// nolint: errcheck
			s.client.ShowMessage(s.lifetime, &protocol.ShowMessageParams{
				Type:    protocol.Error,
				Message: fmt.Sprintf("Failed to connect to Prometheus at %s:\n\n%s ", url, err.Error()),
			})

			s.prometheus = nil

			return err
		}

		// For prometheus version less than 2.14 `api/v1/status/buildinfo` was not supported this can
		// break many function which solely depends on version comparing like `hover`, etc.
		if resp.StatusCode == http.StatusNotFound {
			err = errors.Errorf("Got %d when tried to access %s%s", resp.StatusCode, url, "/api/v1/status/buildinfo")

			// nolint: errcheck
			s.client.ShowMessage(s.lifetime, &protocol.ShowMessageParams{
				Type:    protocol.Error,
				Message: fmt.Sprintf("Failed to access build info via %s:\n\n%s ", url, err.Error()),
			})

			// Setting version 2.13.0 for version less than 2.14 by default.
			s.PrometheusVersion = "2.13.0"

			return err
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		jsonResponse := testResponse{}

		err = json.Unmarshal(bodyBytes, &jsonResponse)

		if err == nil {
			s.PrometheusVersion = jsonResponse.Data.Version

			_, err = semver.Parse(s.PrometheusVersion)
		}

		if err != nil {
			// nolint: errcheck
			s.client.ShowMessage(s.lifetime, &protocol.ShowMessageParams{
				Type:    protocol.Error,
				Message: fmt.Sprintf("Failed to abstract version from Prometheus %s:\n\n%s ", url, err.Error()),
			})

			s.PrometheusVersion = ""

			return err
		}

		resp.Body.Close()
	}

	if err == nil {
		s.PrometheusURL = url
	}

	return err
}

func (s *server) getPrometheus() v1.API {
	s.prometheusMu.Lock()
	defer s.prometheusMu.Unlock()

	if s.prometheus != nil {
		return v1.NewAPI(s.prometheus)
	}

	return nil
}

func (s *server) getPrometheusURL() string {
	s.prometheusMu.Lock()
	defer s.prometheusMu.Unlock()

	return s.PrometheusURL
}

// supportsMetadataAPI checks if language server has access to new metadata APIs.
func (s *server) supportsMetadataAPI() (bool, error) {
	var isCompatible = false // nolint:ineffassign

	requiredVersion := semver.MustParse("2.15.0")
	presentVersion, err := semver.New(s.PrometheusVersion)

	switch {
	case err != nil && !s.headless:
		return isCompatible, err
	case s.headless:
		isCompatible = true
	default:
		isCompatible = presentVersion.GTE(requiredVersion)
	}

	return isCompatible, err
}

// RunTCPServer generates a server listening on the provided TCP Address, creating a new language Server
// instance using plain HTTP for every connection.
func RunTCPServer(ctx context.Context, addr string, config *Config) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go ServerFromStream(ctx, jsonrpc2.NewHeaderStream(conn, conn), config)
	}
}

// StdioServer generates a Server instance talking to stdio.
func StdioServer(ctx context.Context, config *Config) (context.Context, Server) {
	stream := jsonrpc2.NewHeaderStream(os.Stdin, os.Stdout)
	return ServerFromStream(ctx, stream, config)
}
