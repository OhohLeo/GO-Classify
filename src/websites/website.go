package websites

type Website interface {
	GetName() string
	SetConfig(map[string]string) bool
	Search(string) chan Data
}

type Data interface {
	GetType() string
	GetId() string
}
