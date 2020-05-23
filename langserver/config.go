// Copyright 2019 The Prometheus Authors
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
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/lsp/protocol"
	"gopkg.in/yaml.v3"
)

// Config contains the configuration for a server.
type Config struct {
	RPCTrace      string `yaml:"rpc_trace"`
	LogFormat     string `yaml:"log_format"`
	PrometheusURL string `yaml:"prometheus_url"`
	RESTAPIPort   uint64 `yaml:"rest_api_port"`
}

// UnmarshalYAML overrides a function used internally by the yaml.v3 lib.
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tmp := &Config{}
	type plain Config
	if err := unmarshal((*plain)(tmp)); err != nil {
		return err
	}
	if err := tmp.Validate(); err != nil {
		return err
	}
	*c = *tmp
	return nil
}

func (c *Config) unmarshalENV() error {
	prefix := "LANGSERVER"
	conf := &struct {
		RPCTrace      string
		LogFormat     string
		PrometheusURL string
		// the envconfig lib is not able to convert an empty string to the value 0
		// so we have to convert it manually
		RESTAPIPort string
	}{}
	if err := envconfig.Process(prefix, conf); err != nil {
		return err
	}
	if len(conf.RESTAPIPort) > 0 {
		var parseError error
		c.RESTAPIPort, parseError = strconv.ParseUint(conf.RESTAPIPort, 10, 64)
		if parseError != nil {
			return parseError
		}
	}
	c.RPCTrace = conf.RPCTrace
	c.PrometheusURL = conf.PrometheusURL
	c.LogFormat = conf.LogFormat
	return c.Validate()
}

// Validate returns an error if the config is not valid.
func (c *Config) Validate() error {
	if len(c.PrometheusURL) > 0 {
		if _, err := url.Parse(c.PrometheusURL); err != nil {
			return err
		}
	}

	if !regexp.MustCompile("(text|json)?").MatchString(c.LogFormat) {
		return fmt.Errorf(`log Format must be "text", "json" is "%s"`, c.LogFormat)
	}

	return nil
}

// ReadConfig gets the GlobalConfig from a configFile (that is a path to the file).
func ReadConfig(configFile string) (*Config, error) {
	if len(configFile) == 0 {
		fmt.Fprintln(os.Stderr, "No config file provided, configuration is reading from System environment")
		return readConfigFromENV()
	}
	fmt.Fprintln(os.Stderr, "Configuration is reading from configuration file")
	return readConfigFromYAML(configFile)
}

func readConfigFromYAML(configFile string) (*Config, error) {
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	res := new(Config)
	err = yaml.Unmarshal(b, res)
	return res, err
}

func readConfigFromENV() (*Config, error) {
	res := new(Config)
	err := res.unmarshalENV()
	return res, err
}

// DidChangeConfiguration is required by the protocol.Server interface.
func (s *server) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error {
	langserverAddressConfigPath := []string{"promql", "url"}

	if params != nil {
		// nolint: errcheck
		s.client.LogMessage(
			s.lifetime,
			&protocol.LogMessageParams{
				Type:    protocol.Info,
				Message: fmt.Sprintf("Received notification change: %v\n", params),
			})

		setting := params.Settings

		for _, e := range langserverAddressConfigPath {
			m, ok := setting.(map[string]interface{})
			if !ok {
				break
			}

			setting, ok = m[e]
			if !ok {
				break
			}
		}

		if str, ok := setting.(string); ok {
			if err := s.connectPrometheus(str); err != nil {
				// nolint: errcheck
				s.client.LogMessage(ctx, &protocol.LogMessageParams{
					Type:    protocol.Info,
					Message: err.Error(),
				})
			}
		}
	}

	return nil
}
