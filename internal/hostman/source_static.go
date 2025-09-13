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

func (s *StaticSource) GetName() string {
	return s.Name
}

func (s *StaticSource) Validate() error {
	if s.Host == nil && s.Hosts == nil {
		return errors.New("host or hosts must be set")
	}
	if s.Host != nil && s.Hosts != nil {
		return errors.New("only one of host or hosts can be set")
	}
	if s.Host != nil && *s.Host == "" {
		return errors.New("host cannot be empty")
	}

	if s.Ip == "" {
		return errors.New("ip cannot be empty")
	}

	if s.Hosts != nil {
		if len(*s.Hosts) == 0 {
			return errors.New("hosts cannot be empty")
		}
		seen := make(map[string]struct{})
		for _, h := range *s.Hosts {
			if _, ok := seen[h]; ok {
				return errors.New("hosts contains duplicate names")
			}
			seen[h] = struct{}{}
		}
	}

	return nil
}

var _ Source = (*StaticSource)(nil)
