package imports

import (
	"github.com/ohohleo/classify/database"
)

func INIT_DB(db *database.Database) (err error) {

	err = db.AddTable("imports",
		[]string{"id", "name", "ref", "params"})
	if err != nil {
		return
	}

	err = db.AddTable("imports_refs",
		[]string{"id", "name"})
	if err != nil {
		return
	}

	err = db.AddTable("imports_mappings",
		[]string{"imports_id", "collections_id"})
	if err != nil {
		return
	}

	return
}

func INIT_REF_DB(db *database.Database) error {
	return db.InsertRef("imports_refs", REF_IDX2STR)
}

type OnImport func(id uint64, name string, ref Ref, params []byte, collections []string) error

func RetreiveDBImports(db *database.Database, onImport OnImport) (err error) {

	var dbImports []database.GenStruct
	dbImports, err = db.SelectAll("imports")
	if err != nil {
		return
	}

	for _, dbImport := range dbImports {

		var collections []string

		err = db.Select(&collections,
			"SELECT collections.name FROM imports_mappings "+
				"INNER JOIN collections "+
				"WHERE collections.id = imports_mappings.collections_id "+
				"AND imports_mappings.imports_id = ?",
			dbImport.Id)

		if err != nil {
			return
		}

		err = onImport(dbImport.Id, dbImport.Name, Ref(dbImport.Ref),
			dbImport.Params, collections)
	}

	return

}
