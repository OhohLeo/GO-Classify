package exports

type Export interface {
	Check(map[string][]string, []string) error
	Stop()
	GetType() string
	Eq(Export) bool
}
