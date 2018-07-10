package repository

import (
	"env-up-app/backend/types"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type EnvironmentRepository struct {
	environment *types.Environment
	filePath    string
}

func NewEnvironmentRepository(filePath string) (*EnvironmentRepository, error) {
	env, err := loadEnvFromDisk(filePath)
	if err != nil {
		return nil, err
	}
	return &EnvironmentRepository{env, filePath}, nil
}

func (r *EnvironmentRepository) Get() *types.Environment {
	return r.environment
}

func loadEnvFromDisk(filePath string) (*types.Environment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var environment *types.Environment
	err = yaml.NewDecoder(file).Decode(&environment)
	if err != nil {
		return nil, err
	}

	return environment, nil
}
