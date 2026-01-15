package process

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

// Process manages a subprocess.
type Process struct {
	cmd     *exec.Cmd
	started bool
}

// Option configures a Process.
type Option func(*Process)

// WithStdout redirects stdout to the given writer.
func WithStdout(w io.Writer) Option {
	return func(p *Process) {
		if p.cmd != nil {
			p.cmd.Stdout = w
		}
	}
}

// WithStderr redirects stderr to the given writer.
func WithStderr(w io.Writer) Option {
	return func(p *Process) {
		if p.cmd != nil {
			p.cmd.Stderr = w
		}
	}
}

// WithDir sets the working directory.
func WithDir(dir string) Option {
	return func(p *Process) {
		if p.cmd != nil {
			p.cmd.Dir = dir
		}
	}
}

// WithEnv sets environment variables.
func WithEnv(env []string) Option {
	return func(p *Process) {
		if p.cmd != nil {
			p.cmd.Env = env
		}
	}
}

// NewProcess creates a Process for running the given command.
func NewProcess(name string, args ...string) *Process {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return &Process{cmd: cmd}
}

// Apply applies options to the process. Returns the process for chaining.
func (p *Process) Apply(opts ...Option) *Process {
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Start launches the process (non-blocking).
func (p *Process) Start() error {
	if p.started {
		return errors.New("process already started")
	}
	if err := p.cmd.Start(); err != nil {
		return err
	}
	p.started = true
	return nil
}

// Wait waits for the process to exit and returns the exit error if non-zero.
func (p *Process) Wait() error {
	return p.cmd.Wait()
}

// Kill sends SIGKILL to the process.
func (p *Process) Kill() error {
	if p.cmd.Process == nil {
		return nil
	}
	return p.cmd.Process.Kill()
}
