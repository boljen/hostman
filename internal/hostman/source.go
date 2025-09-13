package hostman

type Source interface {
	GetName() string
	GetMapping() (map[string]string, error)
	Validate() error
}
