package configuration

import (
	"gopkg.in/yaml.v3"
	"os"
)

func ReadConfigFromYMLFile(configurationFile string, config *OsekPaHesh) error {
	var f *os.File
	var err error

	if f, err = os.Open(configurationFile); err != nil {
		return err
	}

	defer f.Close()

	return yaml.NewDecoder(f).Decode(config)
}
