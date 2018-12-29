package core

type ImportExportConfig struct {
	General struct {
		Enabled bool `json:"enabled"`
	} `json:"general"`

	Tweak *Tweak `json:"tweak"`
}

func NewImportExportConfig() *ImportExportConfig {
	config := &ImportExportConfig{}

	// Enable import by default
	config.General.Enabled = true

	return config
}

func (i *ImportExportConfig) SetTweak(tweak *Tweak) {
	i.Tweak = tweak
}
