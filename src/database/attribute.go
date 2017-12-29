package database

type AttributeType int

const (
	TEXT AttributeType = iota
	INTEGER
)

var attributeType2str = []string{
	"TEXT",
	"INTEGER",
}

type Attribute struct {
	Type         AttributeType
	IsNotNull    bool
	IsPrimaryKey bool
	IsUnique     bool
}

func (a *Attribute) Create() string {
	res := attributeType2str[a.Type]

	if a.IsNotNull {
		res += " NOT NULL"
	}

	if a.IsPrimaryKey {
		res += " PRIMARY KEY"
	}

	if a.IsUnique {
		res += " UNIQUE"
	}

	return res
}

func (a *Attribute) String() string {
	return attributeType2str[a.Type]
}
