package hostman

import "errors"

type StaticMapping struct {
	Name  string    `hcl:"name,label"`
	Host  *string   `hcl:"host"`
	Hosts *[]string `hcl:"hosts"`
	Ip    string    `hcl:"ip"`
}

func (s *StaticMapping) GetMapping() (map[string]string, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	m := make(map[string]string)
	if s.Host != nil {
		m[*s.Host] = s.Ip
		return m, nil
	} else if s.Hosts != nil {
		for _, h := range *s.Hosts {
			m[h] = s.Ip
		}
		return m, nil
	} else {
		return nil, errors.New("host or hosts must be set")
	}
}

func (s *StaticMapping) GetName() string {
	return s.Name
}

func (s *StaticMapping) Validate() error {
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

var _ Mapping = (*StaticMapping)(nil)
