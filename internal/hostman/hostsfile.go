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

func (f *HostsFile) Write() error {
	return os.WriteFile(f.location, []byte(f.content), 0o600)
}

func (f *HostsFile) Update(project string, filename string, mapping map[string]string) error {
	projectStart := "### hostman-project-start " + project + " ###"
	projectEnd := "### hostman-project-end " + project + " ###"

	projectContent := filename + "\n"
	for domain, ip := range mapping {
		projectContent = projectContent + ip + "\t" + domain + "\n"
	}

	newContent, err := InsertOrReplaceSection(f.content, projectStart, projectEnd, projectContent)
	if err != nil {
		return err
	}

	if err := os.WriteFile(f.location, []byte(newContent), 0o600); err != nil {
		return err
	}

	f.content = newContent

	return nil
}
