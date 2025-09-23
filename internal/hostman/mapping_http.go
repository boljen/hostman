package hostman

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type RemoteHostConfig struct {
	Hosts map[string]string
}

type HTTPMapping struct {
	Name            string `hcl:"name,label"`
	Endpoint        string `hcl:"endpoint"`
	RefreshInterval int    `hcl:"refresh_interval"`
}

func (s *HTTPMapping) GetMapping() (map[string]string, error) {
	cfg, err := s.GetFromRemote()
	if err != nil {
		return nil, err
	}
	if cfg == nil || cfg.Hosts == nil {
		return nil, errors.New("invalid remote config: hosts missing")
	}
	return cfg.Hosts, nil
}

func (s *HTTPMapping) GetName() string {
	return s.Name
}

func (s *HTTPMapping) Validate() error {
	resp, err := http.Get(s.Endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("http error: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var cfg RemoteHostConfig
	if err := json.Unmarshal(body, &cfg); err != nil {
		return err
	}
	if cfg.Hosts == nil {
		return errors.New("invalid response: hosts missing")
	}
	return nil
}

func (s *HTTPMapping) GetFromRemote() (*RemoteHostConfig, error) {
	if s.Endpoint == "" {
		return nil, errors.New("endpoint not set")
	}
	resp, err := http.Get(s.Endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http error: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var cfg RemoteHostConfig
	if err := json.Unmarshal(body, &cfg); err != nil {
		return nil, err
	}
	if cfg.Hosts == nil {
		return nil, errors.New("invalid response: hosts missing")
	}
	return &cfg, nil
}

var _ Mapping = (*HTTPMapping)(nil)
