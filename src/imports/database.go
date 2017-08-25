package imports

import (
	"github.com/ohohleo/classify/database"
)

func INIT_DB(db *database.Database) (err error) {

	err = db.AddTable("imports",
		[]string{"id", "ref", "params"})
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

type OnImport func(id string, ref Ref, params []byte, collections []string) error

func RetreiveDBImports(db *database.Database, onImport OnImport) (err error) {

	// var dbImports []database.GenStruct
	// dbImports, err = db.SelectAll("imports")
	// if err != nil {
	// 	return
	// }

	// for _, dbImport := range dbImports {

	// 	// TODO questionner

	// 	// Add new stored import
	// 	// err = onImport(dbImport.Name, Ref(dbImport.Ref), params, collections)
	// 	// if err != nil {
	// 	// 	return
	// 	// }
	// }

	return

}
