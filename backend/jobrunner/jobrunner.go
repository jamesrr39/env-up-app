package jobrunner

import (
	"env-up-app/backend/types"
	"errors"
	"os/exec"
)

type JobRunner struct{}

func NewJobRunner() *JobRunner {
	return &JobRunner{}
}

var (
	ErrInstanceNotRunning = errors.New("instance not running through env-up")
)

func (r *JobRunner) Run(component *types.Component) error {
	cmd := exec.Command(component.RunCmd)

	err := cmd.Start()
	if err != nil {
		return err
	}

	component.CurrentlyRunningInstance = cmd.Process

	return nil
}

func (r *JobRunner) Stop(component *types.Component) error {
	if component.CurrentlyRunningInstance == nil {
		return ErrInstanceNotRunning
	}

	// err = component.CurrentlyRunningInstance.Kill()
	// if err != nil {
	// 	return err
	// }

	err := component.CurrentlyRunningInstance.Release()
	if err != nil {
		return err
	}

	component.CurrentlyRunningInstance = nil
	return nil
}
