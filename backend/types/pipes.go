package types

type Pipe int

const (
	PipeUnknown Pipe = iota
	PipeStdout
	PipeStderr
)

var pipes = []string{
	"unknown",
	"stdout",
	"stderr",
}

func (p Pipe) String() string {
	return pipes[p]
}
