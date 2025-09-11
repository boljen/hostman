package hostman

type Source interface {
	GetMapping() (map[string]string, error)
}
