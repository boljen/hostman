package hostman

type Mapping interface {
	GetName() string
	GetMapping() (map[string]string, error)
	Validate() error
}
