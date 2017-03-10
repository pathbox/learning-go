package config

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var configDirName = "example"

func GetDefaultConfigDir() (string, error) {
	var configDirLocation string

	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "linux":
		// Use the XDG_CONFIG_HOME variable if it is set, otherwise
		// $HOME/.config/example
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome != "" {
			configDirLocation = xdgConfigHome
		} else {
			configDirLocation = filepath.Join(homeDir, "config", configDirName)
		}

	default:
		hiddenConfigDirname := "." + configDirName
		configDirLocation = filepath.Join(homeDir, hiddenConfigDirname)
	}
	return configDirLocation, nil
}

type Config struct {
	General GeneralOptions
	Keys    map[string]map[string]string
}

type GeneralOptions struct {
	DefaultURLScheme       string
	FormatJSON             bool
	Insecure               bool
	PreserveScrollPosition bool
	Timeout                Duration
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

var DefaultConfig = Config{
	General: GeneralOptions{
		DefaultURLScheme:       "https",
		FormatJSON:             true,
		Insecure:               false,
		PreserveScrollPosition: true,
		Timeout: Duration{
			Duration: 1 * time.Minute,
		},
	},
}

func LoadConfig(configFile string) (*Config, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, errors.New("Config file does not exist.")
	} else if err != nil {
		return nil, err
	}

	conf := DefaultConfig
	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
