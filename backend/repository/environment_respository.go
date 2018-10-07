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
	environment             *types.Environment
	filePath                string
	logMessageChanListeners []chan *types.LogMessage
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
	return &EnvironmentRepository{env, filePath, nil}, nil
}

func (r *EnvironmentRepository) Get() *types.Environment {
	return r.environment
}

func (r *EnvironmentRepository) Start(componentName string) error {
	env := r.Get()
	var component *types.Component
	for _, componentInRange := range env.Components {
		if componentInRange.Name == componentName {
			component = componentInRange
		}
	}

	if component == nil {
		return types.ErrNotFound
	}

	cmd := exec.Command(component.RunCmd)
	cmd.Dir = filepath.Dir(r.filePath)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go r.readFromPipe(component, types.PipeStdout, stdoutPipe)
	go r.readFromPipe(component, types.PipeStderr, stderrPipe)

	err = cmd.Start()
	if err != nil {
		return err
	}
	component.CurrentlyRunningInstance = cmd

	go func() {
		cmd.Wait()
		println("process exited")
	}()

	// return cmd.Wait()
	return nil
}

func (r *EnvironmentRepository) readFromPipe(component *types.Component, pipe types.Pipe, readCloser io.ReadCloser) {
	defer readCloser.Close()

	s := bufio.NewScanner(readCloser)
	for s.Scan() {
		log := types.NewLogMessage(component, pipe, s.Text())
		for _, listener := range r.logMessageChanListeners {
			listener <- log
		}
	}
}

type Listener struct {
	Chan                  chan *types.LogMessage
	environmentRepository *EnvironmentRepository
}

func (l *Listener) Close() {
	for i, thisChan := range l.environmentRepository.logMessageChanListeners {
		if l.Chan == thisChan {
			l.environmentRepository.logMessageChanListeners = append(
				l.environmentRepository.logMessageChanListeners[:i],
				l.environmentRepository.logMessageChanListeners[i+1:]...,
			)
		}
	}
}

func (r *EnvironmentRepository) GetLogMessageChanListener() *Listener {
	listener := &Listener{make(chan *types.LogMessage), r}
	r.logMessageChanListeners = append(r.logMessageChanListeners, listener.Chan)
	return listener
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
