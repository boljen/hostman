package hostman

import (
	"errors"
	"strings"
	"testing"
)

type fakeSource struct {
	name string
	err  error
}

func (f *fakeSource) GetName() string { return f.name }
func (f *fakeSource) GetMapping() (map[string]string, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeSource) Validate() error { return f.err }

var _ Mapping = (*fakeSource)(nil)

func Test_ProjectFileValidate_ReturnsErrorWhenAnySourceFaulty(t *testing.T) {
	invalid := &StaticMapping{Name: "bad", Ip: ""}
	pf := &ProjectFile{Static: []*StaticMapping{invalid}}

	if err := pf.Validate(); err == nil {
		t.Fatalf("expected ProjectFile.Validate to return an error when a source is invalid")
	}
}

func Test_validateIndividualSources_NoSources_OK(t *testing.T) {
	pf := &ProjectFile{}
	if err := pf.validateSourceConfigs(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func Test_validateIndividualSources_CollectsErrors(t *testing.T) {
	// Use StaticMapping with invalid configs to force errors and one valid to pass
	host := "valid.example"
	valid := &StaticMapping{Name: "valid", Host: &host, Ip: "127.0.0.1"}

	hostSet := "some.example"
	invalid1 := &StaticMapping{Name: "bad1", Host: &hostSet, Ip: ""}                                  // ip cannot be empty
	invalid2 := &StaticMapping{Name: "bad2", Host: (*string)(nil), Hosts: &[]string{}, Ip: "1.2.3.4"} // hosts cannot be empty

	pf := &ProjectFile{Static: []*StaticMapping{valid, invalid1, invalid2}}

	err := pf.validateSourceConfigs()
	if err == nil {
		t.Fatalf("expected error from validateSourceConfigs, got nil")
	}
	// The error should mention both invalid cases; we check substrings to avoid depending on join format
	msg := err.Error()
	if !containsAll(msg, []string{"ip cannot be empty", "hosts cannot be empty"}) {
		t.Fatalf("joined error does not contain expected messages: %q", msg)
	}
}

// containsAll checks all substrings are present in s
func containsAll(s string, subs []string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

func TestValidateName(t *testing.T) {

	if err := validateName("domaintest"); err != nil {
		t.Fatalf("validateName failed: %v", err)
	}
}
