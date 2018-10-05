package repository

import (
	"bufio"
	"env-up-app/backend/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type EnvironmentRepository struct {
	environment    *types.Environment
	filePath       string
	LogMessageChan chan *types.LogMessage
}

func NewEnvironmentRepository(filePath string) (*EnvironmentRepository, error) {
	env, err := loadEnvFromDisk(filePath)
	if err != nil {
		return nil, err
	}
	err = env.Validate()
	if err != nil {
		return nil, err
	}
	return &EnvironmentRepository{env, filePath, make(chan *types.LogMessage)}, nil
}

func (r *EnvironmentRepository) Get() *types.Environment {
	return r.environment
}

func (r *EnvironmentRepository) Start(componentName string) error {
	env := r.Get()
	for _, component := range env.Components {
		if component.Name == componentName {
			cmd := exec.Command(component.RunCmd)
			cmd.Dir = filepath.Dir(r.filePath)

			stdoutPipe, err := cmd.StdoutPipe()
			if err != nil {
				return err
			}
			defer stdoutPipe.Close()

			stderrPipe, err := cmd.StderrPipe()
			if err != nil {
				return err
			}
			defer stderrPipe.Close()

			component.CurrentlyRunningInstance = cmd

			go r.readFromPipe(component, types.PipeStdout, stdoutPipe)
			go r.readFromPipe(component, types.PipeStderr, stderrPipe)

			return cmd.Start()
		}
	}

	return types.ErrNotFound
}

func (r *EnvironmentRepository) readFromPipe(component *types.Component, pipe types.Pipe, reader io.Reader) {
	s := bufio.NewScanner(reader)
	for s.Scan() {
		log := types.NewLogMessage(component, pipe, s.Text())
		r.LogMessageChan <- log
	}
}

func (r *EnvironmentRepository) Stop(componentName string) error {

	env := r.Get()
	for _, component := range env.Components {
		if component.Name == componentName {
			component.CurrentlyRunningInstance.Process.Kill()
		}
	}

	return types.ErrNotFound
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
