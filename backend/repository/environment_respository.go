package repository

import (
	"env-up-app/backend/types"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type EnvironmentRepository struct{}

func NewEnvironmentRepository() *EnvironmentRepository {
	return &EnvironmentRepository{}
}

func (r *EnvironmentRepository) Get(filePath string) (*types.Environment, error) {
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
