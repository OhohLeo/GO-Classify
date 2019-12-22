package exports

import (
	"github.com/ohohleo/classify/database"
)

func INIT_DB(db *database.Database) (err error) {

	err = db.AddTable("exports",
		[]string{"id", "name", "ref", "params"})
	if err != nil {
		return
	}

	err = db.AddTable("exports_refs",
		[]string{"id", "name"})
	if err != nil {
		return
	}

	err = db.AddTable("exports_mappings",
		[]string{"exports_id", "collections_id"})
	if err != nil {
		return
	}

	return
}

func INIT_REF_DB(db *database.Database) error {
	return db.InsertRef("exports_refs", REF_IDX2STR)
}

type OnExport func(id uint64, name string, ref Ref, params []byte, collections []string) error

func RetreiveDBExports(db *database.Database, onExport OnExport) (err error) {

	var dbExports []database.GenStruct
	dbExports, err = db.SelectAll("exports")
	if err != nil {
		return
	}

	for _, dbExport := range dbExports {

		var collections []string

		err = db.Select(&collections,
			"SELECT collections.name FROM exports_mappings "+
				"INNER JOIN collections "+
				"WHERE collections.id = exports_mappings.collections_id "+
				"AND exports_mappings.exports_id = ?",
			dbExport.Id)

		if err != nil {
			return
		}

		err = onExport(dbExport.Id, dbExport.Name, Ref(dbExport.Ref),
			dbExport.Params, collections)
	}

	return

}
