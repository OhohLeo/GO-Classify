package websites

const (
	IMDB Type = iota
	TMDB
)

type Type int

func (t Type) String() string {
	return TYPE_IDX2STR[t]
}

var TYPE_IDX2STR = []string{
	"IMDB",
	"TMDB",
}

type Website interface {
	GetName() string
	SetConfig(map[string]string) bool
	Search(string) chan Data
}

type Data interface {
	GetType() string
	GetId() string
}
