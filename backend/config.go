// Steve Phillips / elimisteve
// 2015.11.04

package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/elimisteve/cryptag"
	"github.com/elimisteve/fun"
)

var (
	ErrConfigExists = errors.New("Backend config already exists")
)

type Config struct {
	Name     string
	New      bool `json:"-"`
	Key      *[32]byte
	Local    bool
	DataPath string // Used by backend.FileSystem, other local backends

	Custom map[string]interface{} `json:",omitempty"` // Used by Dropbox, Webserver, other backends
}

func (conf *Config) Save(backendsDir string) error {
	filename := path.Join(backendsDir, conf.Name) + ".json"
	if _, err := os.Stat(filename); err == nil {
		log.Printf("Backend config already exists at %v; NOT overwriting",
			filename)
		return ErrConfigExists
	}

	if err := conf.Canonicalize(); err != nil {
		return err
	}
	b, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, b, 0600)
}

func (conf *Config) Canonicalize() error {
	if conf.Name == "" {
		return errors.New("Storage backend name cannot be empty")
	}
	conf.Name = strings.TrimSuffix(conf.Name, ".json")

	if fun.ContainsAnyStrings(conf.Name, " ", "\t", "\r", "\n") {
		return fmt.Errorf("Storage backend name `%s` contains one or"+
			" more whitespace characters, shouldn't", conf.Name)
	}

	if conf.Key == nil {
		log.Printf("Generating new encryption key for backend `%s`...",
			conf.Name)
		key, err := cryptag.RandomKey()
		if err != nil {
			return err
		}
		conf.Key = key
	}

	if conf.Local && conf.DataPath == "" {
		conf.DataPath = cryptag.LocalDataPath
	}
	conf.DataPath = strings.TrimRight(conf.DataPath, "/\\")

	return nil
}
