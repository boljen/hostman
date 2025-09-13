package hostman

import "os"

type HostsFile struct {
	location string
	content  string
}

func OpenHostsFile(location string) (*HostsFile, error) {
	data, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	file := &HostsFile{
		location: location,
		content:  string(data),
	}
	return file, nil
}
