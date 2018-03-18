package data

import (
	"crypto/md5"
	"encoding/json"
	"math/big"
	//"time"
	//"github.com/pariz/gountries"
)

const (
	GENERIC Ref = iota
	FILE
	MOVIE
	EMAIL
	ATTACHMENT
)

type Ref uint64

func (t Ref) String() string {
	return REF_IDX2STR[t]
}

var REF_IDX2STR = []string{
	"generic",
	"file",
	"movie",
	"email",
	"attachment",
}

var REF_STR2IDX = map[string]Ref{
	REF_IDX2STR[GENERIC]:    GENERIC,
	REF_IDX2STR[FILE]:       FILE,
	REF_IDX2STR[MOVIE]:      MOVIE,
	REF_IDX2STR[EMAIL]:      EMAIL,
	REF_IDX2STR[ATTACHMENT]: ATTACHMENT,
}

type Data interface {
	GetName() string
	// GetDate() time.Time
	// GetCountry() gountries.Country
	GetRef() Ref
}

func GetId(d Data) uint64 {
	res := big.NewInt(0)
	hash := md5.New()
	hash.Write([]byte(d.GetRef().String() + d.GetName()))
	res.SetBytes(hash.Sum(nil))
	return res.Uint64()
}

// Optional link between datas
type HasDependencies interface {
	GetDependencies() []Data
}

// Optional data contents (icons)
type HasContents interface {
	GetContents() map[string]string
}

// Add data functionalities
// - IconsConfig
// - FileConfig
type HasConfig interface {
	NewConfig() Config
}

// Update configuration (IconsConfig, FileConfig...)
type Config interface {
	Update(*json.RawMessage) error
}

// Apply data configuration (FileConfig)
type DataConfig interface {
	ApplyConfig(config Config) error
}

// List of data configs
type Configs map[string]Config

// Handles datas generic interface
func (c *Configs) UnmarshalJSON(src []byte) error {

	datas := make(map[string]*json.RawMessage)

	err := json.Unmarshal(src, &datas)
	if err != nil {
		return err
	}

	for name, config := range *c {
		if rawMsg, ok := datas[name]; ok {
			config.Update(rawMsg)
		}
	}

	return nil
}
