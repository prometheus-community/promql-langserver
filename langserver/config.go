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
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Config contains the configuration for a server
type Config struct {
	RPCTrace      string `yaml:"rpc_trace"`
	PrometheusURL string `yaml:"prometheus_url"`
}

// ParseConfig parses a yaml configuration.
//
// It expects the content of the configuration file as its argument
func ParseConfig(in []byte) (*Config, error) {
	var config Config

	err := yaml.Unmarshal(in, &config)

	return &config, err
}

// ParseConfigFile parses a yaml configuration file.
//
// It expects the path to a configuration file as its argument
func ParseConfigFile(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseConfig(data)
}
