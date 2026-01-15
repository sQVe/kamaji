package process

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewProcess_CreatesCommand(t *testing.T) {
	p := NewProcess("echo", "hello", "world")
	if p == nil {
		t.Fatal("NewProcess() returned nil")
	}
	args := p.cmd.Args
	if len(args) != 3 {
		t.Fatalf("Expected 3 args, got %d: %v", len(args), args)
	}
	if args[0] != "echo" {
		t.Errorf("args[0] = %q, want %q", args[0], "echo")
	}
	if args[1] != "hello" {
		t.Errorf("args[1] = %q, want %q", args[1], "hello")
	}
	if args[2] != "world" {
		t.Errorf("args[2] = %q, want %q", args[2], "world")
	}
}

func TestProcess_Apply(t *testing.T) {
	var buf bytes.Buffer
	p := NewProcess("echo", "test").Apply(WithStdout(&buf))

	if err := p.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := p.Wait(); err != nil {
		t.Fatalf("Wait() error = %v", err)
	}

	output := strings.TrimSpace(buf.String())
	if output != "test" {
		t.Errorf("stdout = %q, want %q", output, "test")
	}
}

func TestProcess_Start(t *testing.T) {
	p := NewProcess("echo", "hello")
	if err := p.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	_ = p.Wait()
}

func TestProcess_Wait_Success(t *testing.T) {
	p := NewProcess("true")
	if err := p.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := p.Wait(); err != nil {
		t.Errorf("Wait() error = %v, want nil for success", err)
	}
}

func TestProcess_Wait_Failure(t *testing.T) {
	p := NewProcess("false")
	if err := p.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := p.Wait(); err == nil {
		t.Error("Wait() error = nil, want error for non-zero exit")
	}
}

func TestProcess_Kill(t *testing.T) {
	p := NewProcess("sleep", "10")
	if err := p.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	time.Sleep(50 * time.Millisecond)

	if err := p.Kill(); err != nil {
		t.Errorf("Kill() error = %v", err)
	}
	_ = p.Wait()
}

func TestProcess_WithStdout(t *testing.T) {
	var buf bytes.Buffer
	p := NewProcess("echo", "hello world").Apply(WithStdout(&buf))

	if err := p.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := p.Wait(); err != nil {
		t.Fatalf("Wait() error = %v", err)
	}

	if output := strings.TrimSpace(buf.String()); output != "hello world" {
		t.Errorf("stdout = %q, want %q", output, "hello world")
	}
}

func TestProcess_WithStderr(t *testing.T) {
	var buf bytes.Buffer
	p := NewProcess("sh", "-c", "echo error >&2").Apply(WithStderr(&buf))

	if err := p.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := p.Wait(); err != nil {
		t.Fatalf("Wait() error = %v", err)
	}

	if output := strings.TrimSpace(buf.String()); output != "error" {
		t.Errorf("stderr = %q, want %q", output, "error")
	}
}

func TestProcess_WithDir(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	p := NewProcess("pwd").Apply(WithStdout(&buf), WithDir(dir))

	if err := p.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := p.Wait(); err != nil {
		t.Fatalf("Wait() error = %v", err)
	}

	if output := strings.TrimSpace(buf.String()); output != dir {
		t.Errorf("pwd = %q, want %q", output, dir)
	}
}

func TestProcess_WithEnv(t *testing.T) {
	var buf bytes.Buffer
	p := NewProcess("sh", "-c", "echo $TEST_VAR").Apply(
		WithStdout(&buf),
		WithEnv([]string{"TEST_VAR=hello"}),
	)

	if err := p.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := p.Wait(); err != nil {
		t.Fatalf("Wait() error = %v", err)
	}

	if output := strings.TrimSpace(buf.String()); output != "hello" {
		t.Errorf("TEST_VAR = %q, want %q", output, "hello")
	}
}
