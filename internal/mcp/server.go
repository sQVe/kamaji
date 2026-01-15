package mcp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sqve/kamaji/internal/version"
)

type Server struct {
	mu           sync.Mutex
	port         int
	actualPort   int
	mcpServer    *server.MCPServer
	httpServer   *http.Server
	started      bool
	signals      chan Signal
	closeSignals sync.Once
}

type Option func(*Server)

// WithPort sets the listening port. Use 0 for dynamic assignment.
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		mcpServer: server.NewMCPServer("kamaji", version.Version,
			server.WithToolCapabilities(true),
		),
		signals: make(chan Signal, 10),
	}
	for _, opt := range opts {
		opt(s)
	}
	s.registerTools()
	return s
}

// Signals returns a receive-only channel for tool call events.
func (s *Server) Signals() <-chan Signal {
	return s.signals
}

func (s *Server) handleTaskComplete(ctx context.Context, req mcp.CallToolRequest, args TaskCompleteArgs) (*mcp.CallToolResult, error) {
	result, err := HandleTaskComplete(ctx, req, args)
	if err != nil {
		return result, err
	}
	if !result.IsError {
		select {
		case s.signals <- Signal{Tool: SignalToolTaskComplete, Status: args.Status, Summary: args.Summary}:
		default:
			slog.Debug("dropped signal: channel full", "tool", SignalToolTaskComplete)
		}
	}
	return result, nil
}

func (s *Server) handleNoteInsight(ctx context.Context, req mcp.CallToolRequest, args NoteInsightArgs) (*mcp.CallToolResult, error) {
	result, err := HandleNoteInsight(ctx, req, args)
	if err != nil {
		return result, err
	}
	if !result.IsError {
		select {
		case s.signals <- Signal{Tool: SignalToolNoteInsight, Summary: args.Text}:
		default:
			slog.Debug("dropped signal: channel full", "tool", SignalToolNoteInsight)
		}
	}
	return result, nil
}

func (s *Server) registerTools() {
	taskCompleteTool := mcp.NewTool("task_complete",
		mcp.WithDescription("Signal task completion"),
		mcp.WithString("status", mcp.Required(), mcp.Description("pass or fail")),
		mcp.WithString("summary", mcp.Required(), mcp.Description("what was done or why it failed")),
	)
	s.mcpServer.AddTool(taskCompleteTool, mcp.NewTypedToolHandler(s.handleTaskComplete))

	noteInsightTool := mcp.NewTool("note_insight",
		mcp.WithDescription("Record discoveries useful for future tasks"),
		mcp.WithString("text", mcp.Required(), mcp.Description("insight to record")),
	)
	s.mcpServer.AddTool(noteInsightTool, mcp.NewTypedToolHandler(s.handleNoteInsight))
}

// Start returns the port once listening.
func (s *Server) Start() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return 0, errors.New("server already started")
	}

	addr := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("failed to listen: %w", err)
	}

	s.actualPort = listener.Addr().(*net.TCPAddr).Port

	mcpHandler := server.NewStreamableHTTPServer(s.mcpServer)
	mux := http.NewServeMux()
	mux.Handle("/mcp", mcpHandler)

	s.httpServer = &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second, // prevent slow-loris attacks
	}
	s.started = true

	// Serve returns ErrServerClosed on graceful shutdown, which we ignore.
	go func() { _ = s.httpServer.Serve(listener) }()

	return s.actualPort, nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started || s.httpServer == nil {
		return nil
	}

	s.started = false
	s.closeSignals.Do(func() { close(s.signals) })
	return s.httpServer.Shutdown(ctx)
}

// Port returns 0 if not started.
func (s *Server) Port() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.actualPort
}
