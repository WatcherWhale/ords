package config

import (
	"os"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

var (
	Config ConfigModel

	defaultConfig ConfigModel = ConfigModel{
		Registries: map[string]string{
			"docker": "index.docker.io",
			"github": "ghcr.io",
		},
	}
)

type ConfigModel struct {
    Registries map[string]string `yaml:"registries"`
}

func ReadConfig() error {
    path, err := xdg.ConfigFile("ords/ords.yaml")

    if err != nil {
        return err
    }

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = writeDefaultConfig(path)

		if err != nil {
			return err
		}
	}

	bytes, err := os.ReadFile(path)

    if err != nil {
        return err
    }

	err = yaml.Unmarshal(bytes, &Config)

    if err != nil {
        return err
    }

    return nil
}

func writeDefaultConfig(path string) error {
	bytes, err := yaml.Marshal(&defaultConfig)

	if err != nil {
		return err
	}

	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(bytes)

	if err != nil {
		return err
	}

	return nil
}
