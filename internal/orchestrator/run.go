package orchestrator

import (
	"context"
	"errors"
	"os"

	"github.com/sqve/kamaji/internal/config"
	"github.com/sqve/kamaji/internal/domain"
	"github.com/sqve/kamaji/internal/git"
	"github.com/sqve/kamaji/internal/mcp"
	"github.com/sqve/kamaji/internal/output"
	"github.com/sqve/kamaji/internal/process"
	"github.com/sqve/kamaji/internal/prompt"
	"github.com/sqve/kamaji/internal/statemachine"
)

// ProcessSpawner abstracts process creation for testing.
type ProcessSpawner interface {
	Spawn(cfg process.SpawnConfig) (*process.SpawnResult, error)
}

// defaultSpawner uses the real process.SpawnClaude.
type defaultSpawner struct{}

func (defaultSpawner) Spawn(cfg process.SpawnConfig) (*process.SpawnResult, error) {
	return process.SpawnClaude(cfg)
}

// commandSpawner runs an arbitrary command instead of Claude.
type commandSpawner struct {
	cmd string
}

func (s commandSpawner) Spawn(cfg process.SpawnConfig) (*process.SpawnResult, error) {
	return process.SpawnCommand(s.cmd, cfg)
}

// RunConfig configures the Run function.
type RunConfig struct {
	WorkDir    string         // Required: project directory
	SprintPath string         // Required: path to kamaji.yaml
	Spawner    ProcessSpawner // Optional: defaults to real spawner
	SpawnerCmd string         // Optional: override spawner with command
}

// RunResult contains the outcome of a sprint execution.
type RunResult struct {
	Success     bool   // true if sprint completed
	TasksRun    int    // count of tasks executed
	Stuck       bool   // true if stuck threshold hit
	StuckReason string // last failure summary if stuck
}

// Run executes a sprint, processing tasks sequentially until completion or stuck.
func Run(ctx context.Context, cfg RunConfig) (*RunResult, error) {
	if cfg.WorkDir == "" {
		return nil, errors.New("WorkDir is required")
	}
	if cfg.SprintPath == "" {
		return nil, errors.New("SprintPath is required")
	}

	sprint, err := config.LoadSprint(cfg.SprintPath)
	if err != nil {
		return nil, err
	}

	state, err := config.LoadState(cfg.WorkDir)
	if err != nil {
		return nil, err
	}

	if len(sprint.Tickets) == 0 {
		output.PrintInfo("Sprint has no tasks")
		return &RunResult{Success: true}, nil
	}

	server := mcp.NewServer(mcp.WithPort(0))
	port, err := server.Start()
	if err != nil {
		return nil, err
	}
	//nolint:contextcheck // Fresh context needed since caller context may be cancelled
	defer func() { _ = server.Shutdown(context.Background()) }()

	handler := NewHandler(cfg.WorkDir, state, sprint)

	spawner := cfg.Spawner
	if spawner == nil {
		if cfg.SpawnerCmd != "" {
			spawner = commandSpawner{cmd: cfg.SpawnerCmd}
		} else {
			spawner = defaultSpawner{}
		}
	}

	var tasksRun int

	for {
		select {
		case <-ctx.Done():
			return &RunResult{TasksRun: tasksRun}, ctx.Err()
		default:
		}

		taskInfo := statemachine.NextTask(state, sprint)
		if taskInfo == nil {
			output.PrintSprintComplete(sprint, state)
			return &RunResult{Success: true, TasksRun: tasksRun}, nil
		}

		// Create branch only at the start of a new ticket. CurrentTask==0 means we're
		// on the first task, and FailureCount==0 means this is not a retry. On retry,
		// the branch already exists from the initial attempt. When advancing to a new
		// ticket, statemachine resets both counters, triggering branch creation again.
		if state.CurrentTask == 0 && state.FailureCount == 0 {
			if err := createTicketBranch(cfg.WorkDir, sprint.BaseBranch, taskInfo.Ticket); err != nil {
				return nil, err
			}
		}

		output.PrintTaskStart(taskInfo, sprint)

		result, err := runTask(ctx, &taskContext{
			cfg:      cfg,
			spawner:  spawner,
			sprint:   sprint,
			state:    state,
			port:     port,
			taskInfo: taskInfo,
			server:   server,
		})
		if err != nil {
			return &RunResult{TasksRun: tasksRun}, err
		}

		tasksRun++

		output.PrintSignal(mcp.Signal{
			Tool:    mcp.SignalToolTaskComplete,
			Status:  result.Status,
			Summary: result.Summary,
		})

		if result.Passed() {
			if err := handler.OnPass(taskInfo.Ticket.Name, taskInfo.Task.Description, result.Summary); err != nil {
				return nil, err
			}
		} else {
			if err := handler.OnFail(taskInfo.Ticket.Name, taskInfo.Task.Description, result.Summary); err != nil {
				return nil, err
			}

			if handler.IsStuck() {
				if err := handler.OnStuck(); err != nil {
					return nil, err
				}
				return &RunResult{TasksRun: tasksRun, Stuck: true, StuckReason: result.Summary}, nil
			}
		}
	}
}

func createTicketBranch(workDir, baseBranch string, ticket *domain.Ticket) error {
	output.PrintTicketStart(ticket)

	err := git.CreateBranch(workDir, baseBranch, ticket.Branch)
	if err != nil {
		if errors.Is(err, git.ErrBranchExists) {
			output.PrintInfo("Using existing branch: " + ticket.Branch)
			return nil
		}
		return err
	}

	output.PrintBranchCreated(ticket.Branch)
	return nil
}

// taskContext groups parameters needed for task execution.
type taskContext struct {
	cfg      RunConfig
	spawner  ProcessSpawner
	sprint   *domain.Sprint
	state    *domain.State
	port     int
	taskInfo *statemachine.TaskInfo
	server   *mcp.Server
}

func runTask(ctx context.Context, tc *taskContext) (TaskResult, error) {
	promptText, err := prompt.AssembleContext(tc.sprint, tc.state, tc.cfg.WorkDir)
	if err != nil {
		return TaskResult{}, err
	}

	spawnResult, err := tc.spawner.Spawn(process.SpawnConfig{
		Prompt:  promptText,
		MCPPort: tc.port,
		WorkDir: tc.cfg.WorkDir,
		Stdout:  output.NewInfoWriter(os.Stdout),
		Stderr:  output.NewErrorWriter(os.Stderr),
	})
	if err != nil {
		return TaskResult{}, err
	}
	defer func() {
		if spawnResult.ConfigPath != "" {
			_ = os.Remove(spawnResult.ConfigPath)
		}
	}()

	done := make(chan struct{})
	go func() {
		_ = spawnResult.Process.Wait()
		close(done)
	}()

	ticketName := tc.taskInfo.Ticket.Name
	for {
		select {
		case <-ctx.Done():
			_ = spawnResult.Process.Kill()
			return TaskResult{}, ctx.Err()
		case sig, ok := <-tc.server.Signals():
			if !ok {
				_ = spawnResult.Process.Kill()
				<-done
				return NoSignalResult(), nil
			}
			if sig.Tool == mcp.SignalToolNoteInsight {
				_ = config.RecordInsight(tc.cfg.WorkDir, ticketName, sig.Summary)
				output.PrintSignal(sig)
				continue
			}
			select {
			case <-done:
			case <-ctx.Done():
				_ = spawnResult.Process.Kill()
				<-done
				return TaskResult{}, ctx.Err()
			}
			return ResultFromSignal(sig), nil
		case <-done:
			// Process exited. Drain any pending signals to capture insights and
			// task completion that arrived concurrently with process exit.
			for {
				select {
				case sig, ok := <-tc.server.Signals():
					if !ok {
						return NoSignalResult(), nil
					}
					if sig.Tool == mcp.SignalToolNoteInsight {
						_ = config.RecordInsight(tc.cfg.WorkDir, ticketName, sig.Summary)
						output.PrintSignal(sig)
						continue
					}
					return ResultFromSignal(sig), nil
				default:
					return NoSignalResult(), nil
				}
			}
		}
	}
}
