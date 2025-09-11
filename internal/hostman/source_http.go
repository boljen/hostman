package hostman

import "errors"

type HTTPSource struct {
	Name            string `hcl:"name,label"`
	Endpoint        string `hcl:"endpoint"`
	RefreshInterval int    `hcl:"refresh_interval"`
}

func (s *HTTPSource) GetMapping() (map[string]string, error) {
	return nil, errors.New("not yet implemented")
}

var _ Source = (*HTTPSource)(nil)
