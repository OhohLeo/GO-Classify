package core

type CollectionImport struct {
	Collection *Collection
	Config     *ImportExportConfig
}

func NewCollectionImport(collection *Collection) *CollectionImport {
	return &CollectionImport{
		Collection: collection,
		Config:     NewImportExportConfig(),
	}
}
