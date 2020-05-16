package langserver

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalENV(t *testing.T) {
	testSuites := []struct {
		title     string
		variables map[string]string
		expected  *Config
	}{
		{
			title:     "empty config",
			variables: map[string]string{},
			expected:  &Config{},
		},
		{
			title: "full config",
			variables: map[string]string{
				"LANGSERVER_RPCTRACE":      "text",
				"LANGSERVER_PROMETHEUSURL": "http://localhost:9090",
				"LANGSERVER_RESTAPIPORT":   "8080",
			},
			expected: &Config{
				RPCTrace:      "text",
				PrometheusURL: "http://localhost:9090",
				RESTAPIPort:   8080,
			},
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.title, func(t *testing.T) {
			os.Clearenv()
			for k, v := range testSuite.variables {
				os.Setenv(k, v)
			}
			conf, err := ReadConfig("")
			assert.Nil(t, err)
			assert.Equal(t, testSuite.expected, conf)
		})
	}

}
