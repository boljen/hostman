package hostman

import "errors"

type StaticSource struct {
	Name  string    `hcl:"name,label"`
	Host  *string   `hcl:"host"`
	Hosts *[]string `hcl:"hosts"`
	Ip    string    `hcl:"ip"`
}

func (s *StaticSource) GetMapping() (map[string]string, error) {
	return nil, errors.New("not yet implemented")
}

var _ Source = (*StaticSource)(nil)
