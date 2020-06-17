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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	kitlog "github.com/go-kit/kit/log"
	promClient "github.com/prometheus-community/promql-langserver/prometheus"

	"github.com/prometheus-community/promql-langserver/langserver"
	"github.com/prometheus-community/promql-langserver/rest"
)

func main() {
	configFilePath := flag.String("config-file", "", "Configuration file for the language server. If unset, the configuration will be retrieved from environment")

	flag.Parse()

	config, err := langserver.ReadConfig(*configFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading config file:", err.Error())
		os.Exit(1)
	}
	if config.RESTAPIPort != 0 {
		fmt.Fprintln(os.Stderr, "REST API: Listening on port ", config.RESTAPIPort)
		prometheusClient, err := promClient.NewClient(config.PrometheusURL)
		if err != nil {
			log.Fatal(err)
		}

		var logger kitlog.Logger

		switch config.LogFormat {
		case langserver.JSONFormat:
			logger = kitlog.NewJSONLogger(os.Stderr)
		case langserver.TextFormat:
			logger = kitlog.NewLogfmtLogger(os.Stderr)
		default:
			log.Fatalf(`invalid log format: "%s"`, config.LogFormat)
		}

		logger = kitlog.NewSyncLogger(logger)

		handler, err := rest.CreateInstHandler(context.Background(), prometheusClient, logger)
		if err != nil {
			log.Fatal(err)
		}

		err = http.ListenAndServe(fmt.Sprint(":", config.RESTAPIPort), handler)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, s := langserver.StdioServer(context.Background(), config)
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
