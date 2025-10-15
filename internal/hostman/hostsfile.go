package hostman

import (
	"os"
	"regexp"
)

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

	projectContent := "### location: " + filename + "\n"
	for domain, ip := range mapping {
		projectContent = projectContent + ip + "\t" + domain + "\n"
	}

	newContent, err := InsertOrReplaceSection(f.content, projectStart, projectEnd, projectContent)
	if err != nil {
		return err
	}

	f.content = newContent

	return nil
}

func (f *HostsFile) Save() error {
	return os.WriteFile(f.location, []byte(f.content), 0o600)
}

var projectNameRegexp = regexp.MustCompile(`### hostman-project-start\s+(.*?)\s+###`)

func (f *HostsFile) GetProjects() []string {
	matches := projectNameRegexp.FindAllStringSubmatch(f.content, -1)
	projects := make([]string, 0, len(matches))
	for i := 0; i < len(matches); i++ {
		projects = append(projects, matches[i][1])
	}
	return projects
}

func (f *HostsFile) GetProjectData(project string) (string, error) {
	projectStart := "### hostman-project-start " + project + " ###"
	projectEnd := "### hostman-project-end " + project + " ###"
	return GetSection(f.content, projectStart, projectEnd)
}

func (f *HostsFile) RemoveAllProjects() error {
	projects := f.GetProjects()
	for _, p := range projects {
		if err := f.RemoveProject(p); err != nil {
			return err
		}
	}
	return nil
}

func (f *HostsFile) RemoveProject(project string) error {
	projectStart := "### hostman-project-start " + project + " ###"
	projectEnd := "### hostman-project-end " + project + " ###"

	newContent, err := RemoveSection(f.content, projectStart, projectEnd)
	if err != nil {
		return err
	}
	f.content = newContent

	return nil
}
