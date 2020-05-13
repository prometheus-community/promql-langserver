package prometheus

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/blang/semver"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
)

var (
	requiredVersion = semver.MustParse("2.15.0")
)

func buildGenericRoundTripper(connectionTimeout time.Duration) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   connectionTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 30 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, // nolint: gas, gosec
	}
}

func buildStatusRequest(prometheusURL string) (*http.Request, error) {
	finalURL, err := url.Parse(prometheusURL)
	if err != nil {
		return nil, err
	}
	finalURL.Path = "/api/v1/status/buildinfo"
	httpRequest, err := http.NewRequest(http.MethodGet, finalURL.String(), nil)
	if err != nil {
		return nil, err
	}
	// set the accept content type
	httpRequest.Header.Set("Accept", "application/json")
	return httpRequest, nil
}

type buildInfoResponse struct {
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

type Client interface {
	Metadata(metric string) (v1.Metadata, error)
	AllMetadata() (map[string][]v1.Metadata, error)
	ChangeDataSource(prometheusURL string) error
}

// httpClient is an implementation of the interface Client.
// You should use this instance directly and not the other one (compatibleHTTPClient and notCompatibleHTTPClient)
// because it will manage which sub instance of the Client to use (like a factory)
type httpClient struct {
	Client
	requestTimeout time.Duration
	mutex          sync.RWMutex
	subClient      Client
}

func NewClient(prometheusURL string, requestTimeout time.Duration) (Client, error) {
	c := &httpClient{
		requestTimeout: requestTimeout,
	}
	if err := c.ChangeDataSource(prometheusURL); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *httpClient) Metadata(metric string) (v1.Metadata, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.subClient.Metadata(metric)
}

func (c *httpClient) AllMetadata() (map[string][]v1.Metadata, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.subClient.AllMetadata()
}

func (c *httpClient) ChangeDataSource(prometheusURL string) error {
	prometheusHTTPClient, err := api.NewClient(api.Config{
		RoundTripper: buildGenericRoundTripper(c.requestTimeout * time.Second),
		Address:      prometheusURL,
	})
	if err != nil {
		return err
	}

	isCompatible, err := c.isCompatible(prometheusURL)
	if err != nil {
		return err
	}

	// only lock when we are sure we are going to change the instance of the sub client
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if isCompatible {
		c.subClient = &compatibleHTTPClient{
			requestTimeout:   c.requestTimeout,
			prometheusClient: v1.NewAPI(prometheusHTTPClient),
		}
	} else {
		c.subClient = &notCompatibleHTTPClient{
			requestTimeout:   c.requestTimeout,
			prometheusClient: v1.NewAPI(prometheusHTTPClient),
		}
	}

	return nil
}

func (c *httpClient) isCompatible(prometheusURL string) (bool, error) {
	httpRequest, err := buildStatusRequest(prometheusURL)
	if err != nil {
		return false, err
	}
	httpClient := &http.Client{
		Transport: buildGenericRoundTripper(c.requestTimeout * time.Second),
		Timeout:   c.requestTimeout * time.Second,
	}
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return false, err
	}

	// For prometheus version less than 2.14 `api/v1/status/buildinfo` was not supported this can
	// break many function which solely depends on version comparing like `hover`, etc.
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	defer resp.Body.Close() // nolint: errcheck
	if resp.Body != nil {
		data, err := ioutil.ReadAll(resp.Body)
		jsonResponse := buildInfoResponse{}
		err = json.Unmarshal(data, &jsonResponse)
		if err != nil {
			return false, err
		}
		currentVersion, err := semver.New(jsonResponse.Data.Version)
		if err != nil {
			return false, err
		}
		return currentVersion.GTE(requiredVersion), nil
	}
	return false, nil
}

// compatibleHTTPClient must be used to contact a distant prometheus with a version >= v2.15
type compatibleHTTPClient struct {
	Client
	requestTimeout   time.Duration
	prometheusClient v1.API
}

func (c *compatibleHTTPClient) Metadata(metric string) (v1.Metadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout*time.Second)
	defer cancel()

	metadata, err := c.prometheusClient.Metadata(ctx, metric, "1")
	if err != nil {
		return v1.Metadata{}, err
	}
	if len(metadata) <= 0 {
		return v1.Metadata{}, nil
	}
	return v1.Metadata{
		Type: metadata[metric][0].Type,
		Help: metadata[metric][0].Help,
		Unit: metadata[metric][0].Unit,
	}, nil
}

func (c *compatibleHTTPClient) AllMetadata() (map[string][]v1.Metadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout*time.Second)
	defer cancel()
	return c.prometheusClient.Metadata(ctx, "", "")
}

func (c *compatibleHTTPClient) ChangeDataSource(_ string) error {
	return fmt.Errorf("method not supported")
}

// notCompatibleHTTPClient must be used to contact a distant prometheus with a version < v2.15
type notCompatibleHTTPClient struct {
	Client
	requestTimeout   time.Duration
	prometheusClient v1.API
}

func (c *notCompatibleHTTPClient) Metadata(metric string) (v1.Metadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout*time.Second)
	defer cancel()

	metadata, err := c.prometheusClient.TargetsMetadata(ctx, "", metric, "1")
	if err != nil {
		return v1.Metadata{}, err
	}
	if len(metadata) <= 0 {
		return v1.Metadata{}, nil
	}
	return v1.Metadata{
		Type: metadata[0].Type,
		Help: metadata[0].Help,
		Unit: metadata[0].Unit,
	}, nil
}

func (c *notCompatibleHTTPClient) AllMetadata() (map[string][]v1.Metadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout*time.Second)
	defer cancel()

	metricNames, _, err := c.prometheusClient.LabelValues(ctx, "__name__")
	if err != nil {
		return nil, err
	}
	allMetadata := make(map[string][]v1.Metadata)
	for _, name := range metricNames {
		allMetadata[string(name)] = []v1.Metadata{{}}
	}
	return allMetadata, nil
}

func (c *notCompatibleHTTPClient) ChangeDataSource(_ string) error {
	return fmt.Errorf("method not supported")
}
