package repository

import (
	"env-up-app/backend/types"
	"fmt"
	"os"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

type ConfigRepository struct {
	configFilePath string
	config         *types.Configuration
	mu             *sync.Mutex
}

func NewConfigRepository(configFilePath string) (*ConfigRepository, error) {
	config, err := getFromDisk(configFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			// if the error isn't that the file doesn't exist, error out
			return nil, err
		}
		config = &types.Configuration{}
	}

	if len(config.EnvironmentsFilePaths) == 0 {
		config.EnvironmentsFilePaths = []string{}
	}

	return &ConfigRepository{configFilePath, config, &sync.Mutex{}}, nil
}

func getFromDisk(configFilePath string) (*types.Configuration, error) {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config *types.Configuration
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (r *ConfigRepository) Get() *types.Configuration {
	return r.config
}

func (r *ConfigRepository) AddEnvironmentPath(environmentPath string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	copyOfFilePaths := make([]string, len(r.config.EnvironmentsFilePaths))
	copy(copyOfFilePaths, r.config.EnvironmentsFilePaths)

	err := r.set(&types.Configuration{EnvironmentsFilePaths: copyOfFilePaths})
	if err != nil {
		return err
	}

	r.config.EnvironmentsFilePaths = copyOfFilePaths
	return nil
}

func (r *ConfigRepository) set(config *types.Configuration) error {
	configFile, err := os.Open(r.configFilePath)
	if err != nil {
		return fmt.Errorf("failed to open config file at '%s'. Error: '%s'", r.configFilePath, err)
	}
	defer configFile.Close()

	err = yaml.NewEncoder(configFile).Encode(config)
	if err != nil {
		return fmt.Errorf("failed to read config file at '%s'. Error: '%s'", r.configFilePath, err)
	}

	r.config = config

	return nil
}
