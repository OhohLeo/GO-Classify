package collections

import (
	"github.com/ohohleo/classify/database"
)

func INIT_DB(db *database.Database) (err error) {

	err = db.AddTable("collections",
		[]string{"id", "name", "ref", "params"})
	if err != nil {
		return
	}

	err = db.AddTable("collections_refs",
		[]string{"id", "name"})
	if err != nil {
		return
	}

	return
}

func INIT_REF_DB(db *database.Database) error {
	return db.InsertRef("collections_refs", REF_IDX2STR)
}

type OnCollection func(id uint64, name string, ref Ref, params []byte) error

func RetreiveDBCollections(db *database.Database, onCollection OnCollection) (err error) {

	var dbCollections []database.GenStruct
	dbCollections, err = db.SelectAll("collections")
	if err != nil {
		return
	}

	for _, dbCollection := range dbCollections {

		// Get collection params
		var params []byte
		err = dbCollection.GetParams(&params)
		if err != nil {
			return
		}

		// Add new stored collection
		err = onCollection(dbCollection.Id, dbCollection.Name, Ref(dbCollection.Ref), params)
		if err != nil {
			return
		}
	}

	return
}
