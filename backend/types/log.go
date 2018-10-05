package types

type LogMessage struct {
	ComponentName string `json:"componentName"`
	Pipe          Pipe   `json:"pipe"`
	Message       string `json:"message"`
}

func NewLogMessage(component *Component, pipe Pipe, message string) *LogMessage {
	return &LogMessage{component.Name, pipe, message}
}
