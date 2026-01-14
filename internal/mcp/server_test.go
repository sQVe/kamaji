package mcp

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestNewServer_DefaultPort(t *testing.T) {
	s := NewServer()
	if s.port != 0 {
		t.Errorf("NewServer() default port = %d, want 0", s.port)
	}
}

func TestNewServer_WithPort(t *testing.T) {
	s := NewServer(WithPort(8080))
	if s.port != 8080 {
		t.Errorf("NewServer(WithPort(8080)) port = %d, want 8080", s.port)
	}
}

func TestServer_Start(t *testing.T) {
	s := NewServer()

	port, err := s.Start()
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	defer func() { _ = s.Shutdown(context.Background()) }()

	if port == 0 {
		t.Error("Start() returned port 0, want non-zero")
	}

	// Verify server is listening by making a request
	resp, err := http.Get("http://localhost:" + strconv.Itoa(port) + "/mcp")
	if err != nil {
		t.Fatalf("HTTP GET error = %v", err)
	}
	_ = resp.Body.Close()
}

func TestServer_StartTwice(t *testing.T) {
	s := NewServer()

	_, err := s.Start()
	if err != nil {
		t.Fatalf("Start() first call error = %v", err)
	}
	defer func() { _ = s.Shutdown(context.Background()) }()

	_, err = s.Start()
	if err == nil {
		t.Error("Start() second call error = nil, want error")
	}
}

func TestServer_Port(t *testing.T) {
	s := NewServer()

	// Before Start, Port returns 0
	if got := s.Port(); got != 0 {
		t.Errorf("Port() before Start = %d, want 0", got)
	}

	port, err := s.Start()
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	defer func() { _ = s.Shutdown(context.Background()) }()

	// After Start, Port returns the actual port
	if got := s.Port(); got != port {
		t.Errorf("Port() after Start = %d, want %d", got, port)
	}
}

func TestServer_Shutdown(t *testing.T) {
	s := NewServer()

	_, err := s.Start()
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	port := s.Port()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}

	// Give the server a moment to fully shut down
	time.Sleep(50 * time.Millisecond)

	// Verify server is no longer listening
	resp, err := http.Get("http://localhost:" + strconv.Itoa(port) + "/mcp")
	if err == nil {
		_ = resp.Body.Close()
		t.Error("Server still responding after Shutdown")
	}
}

func TestServer_ToolsRegistered(t *testing.T) {
	s := NewServer()

	// Create in-process client to verify tools
	c, err := client.NewInProcessClient(s.mcpServer)
	if err != nil {
		t.Fatalf("NewInProcessClient() error = %v", err)
	}
	defer func() { _ = c.Close() }()

	ctx := context.Background()
	_, err = c.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: "2024-11-05",
			ClientInfo: mcp.Implementation{
				Name:    "test-client",
				Version: "1.0.0",
			},
		},
	})
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	tools, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		t.Fatalf("ListTools() error = %v", err)
	}

	// Verify both tools are registered
	toolNames := make(map[string]bool)
	for _, tool := range tools.Tools {
		toolNames[tool.Name] = true
	}

	if !toolNames["task_complete"] {
		t.Error("tool task_complete not registered")
	}
	if !toolNames["note_insight"] {
		t.Error("tool note_insight not registered")
	}
}
