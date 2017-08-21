package websites

import (
	"github.com/ohohleo/classify/database"
)

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

var DB_DETAILS database.Table = database.Table{
	Name: "websites",
	Attributes: map[string]*database.Attribute{
		"id": &database.Attribute{
			Type:         database.INTEGER,
			IsPrimaryKey: true,
		},
		"collection_id": &database.Attribute{
			Type: database.INTEGER,
		},
		"website_ref_id": &database.Attribute{
			Type: database.INTEGER,
		},
	},
}

var DB_TYPES_REF database.Table = database.Table{
	Name: "website_refs",
	Attributes: map[string]*database.Attribute{
		"id": &database.Attribute{
			Type:         database.INTEGER,
			IsPrimaryKey: true,
		},
		"name": &database.Attribute{
			Type:     database.TEXT,
			IsUnique: true,
		},
	},
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
