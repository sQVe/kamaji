package output

import (
	"strings"
	"testing"

	"github.com/sqve/kamaji/internal/config"
	"github.com/sqve/kamaji/internal/domain"
	"github.com/sqve/kamaji/internal/statemachine"
	"github.com/sqve/kamaji/internal/testutil"
)

func TestTaskProgress_PlainMode(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{
				Name:   "login-form",
				Branch: "feat/login-form",
				Tasks: []domain.Task{
					{Description: "Create login component"},
					{Description: "Add validation"},
					{Description: "Add tests"},
				},
			},
			{
				Name:   "signup-form",
				Branch: "feat/signup-form",
				Tasks: []domain.Task{
					{Description: "Create signup component"},
				},
			},
			{
				Name:   "dashboard",
				Branch: "feat/dashboard",
				Tasks: []domain.Task{
					{Description: "Build dashboard"},
				},
			},
		},
	}

	tests := []struct {
		name        string
		ticketIndex int
		taskIndex   int
		expected    string
	}{
		{
			name:        "first ticket first task",
			ticketIndex: 0,
			taskIndex:   0,
			expected:    "[Ticket 1/3] [Task 1/3] Create login component",
		},
		{
			name:        "first ticket second task",
			ticketIndex: 0,
			taskIndex:   1,
			expected:    "[Ticket 1/3] [Task 2/3] Add validation",
		},
		{
			name:        "second ticket first task",
			ticketIndex: 1,
			taskIndex:   0,
			expected:    "[Ticket 2/3] [Task 1/1] Create signup component",
		},
		{
			name:        "last ticket last task",
			ticketIndex: 2,
			taskIndex:   0,
			expected:    "[Ticket 3/3] [Task 1/1] Build dashboard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &statemachine.TaskInfo{
				TicketIndex: tt.ticketIndex,
				TaskIndex:   tt.taskIndex,
				Ticket:      &sprint.Tickets[tt.ticketIndex],
				Task:        &sprint.Tickets[tt.ticketIndex].Tasks[tt.taskIndex],
			}

			got := TaskProgress(info, sprint)
			if got != tt.expected {
				t.Errorf("TaskProgress() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestTaskProgress_StyledMode(t *testing.T) {
	config.SetPlain(false)
	defer config.ResetPlain()

	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{
				Name:   "login-form",
				Branch: "feat/login-form",
				Tasks: []domain.Task{
					{Description: "Create login component"},
					{Description: "Add validation"},
				},
			},
			{
				Name:   "signup-form",
				Branch: "feat/signup-form",
				Tasks: []domain.Task{
					{Description: "Create signup component"},
				},
			},
		},
	}

	info := &statemachine.TaskInfo{
		TicketIndex: 0,
		TaskIndex:   1,
		Ticket:      &sprint.Tickets[0],
		Task:        &sprint.Tickets[0].Tasks[1],
	}

	got := TaskProgress(info, sprint)

	mustContain := []string{"Ticket 1/2", "Task 2/2", "Add validation", ">"}
	for _, want := range mustContain {
		if !strings.Contains(got, want) {
			t.Errorf("TaskProgress() = %q, missing %q", got, want)
		}
	}
}

func TestTaskProgress_SingleTicketSingleTask(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{
				Name:   "only-ticket",
				Branch: "feat/only",
				Tasks: []domain.Task{
					{Description: "Only task"},
				},
			},
		},
	}

	info := &statemachine.TaskInfo{
		TicketIndex: 0,
		TaskIndex:   0,
		Ticket:      &sprint.Tickets[0],
		Task:        &sprint.Tickets[0].Tasks[0],
	}

	got := TaskProgress(info, sprint)
	expected := "[Ticket 1/1] [Task 1/1] Only task"
	if got != expected {
		t.Errorf("TaskProgress() = %q, want %q", got, expected)
	}
}

func TestPrintTicketStart(t *testing.T) {
	ticket := &domain.Ticket{
		Name:   "login-form",
		Branch: "feat/login-form",
	}

	t.Run("plain mode", func(t *testing.T) {
		config.SetPlain(true)
		defer config.ResetPlain()
		output := testutil.CaptureStdout(t, func() {
			PrintTicketStart(ticket)
		})
		testutil.AssertContains(t, output, "Starting ticket: login-form")
		testutil.AssertContains(t, output, "feat/login-form")
	})

	t.Run("styled mode", func(t *testing.T) {
		config.SetPlain(false)
		defer config.ResetPlain()
		output := testutil.CaptureStdout(t, func() {
			PrintTicketStart(ticket)
		})
		testutil.AssertContains(t, output, "Starting ticket: login-form")
		testutil.AssertContains(t, output, "feat/login-form")
	})
}

func TestGitFeedback(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	t.Run("PrintBranchCreated", func(t *testing.T) {
		output := testutil.CaptureStdout(t, func() {
			PrintBranchCreated("feat/login-form")
		})
		testutil.AssertContains(t, output, "Created branch: feat/login-form")
	})

	t.Run("PrintCommitCreated short message", func(t *testing.T) {
		output := testutil.CaptureStdout(t, func() {
			PrintCommitCreated("Add login form")
		})
		testutil.AssertContains(t, output, "Committed: Add login form")
	})

	t.Run("PrintCommitCreated long message truncates", func(t *testing.T) {
		output := testutil.CaptureStdout(t, func() {
			PrintCommitCreated("This is a very long commit message that exceeds fifty characters")
		})
		testutil.AssertContains(t, output, "Committed:")
		testutil.AssertContains(t, output, "...")
	})

	t.Run("PrintResetPerformed", func(t *testing.T) {
		output := testutil.CaptureStdout(t, func() {
			PrintResetPerformed()
		})
		testutil.AssertContains(t, output, "Reset to HEAD")
	})
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "short string unchanged",
			input:    "short",
			maxLen:   10,
			expected: "short",
		},
		{
			name:     "exact length unchanged",
			input:    "exact",
			maxLen:   5,
			expected: "exact",
		},
		{
			name:     "long string truncated",
			input:    "this is a long string",
			maxLen:   10,
			expected: "this is...",
		},
		{
			name:     "fifty char truncation",
			input:    "This is a very long commit message that exceeds fifty characters",
			maxLen:   50,
			expected: "This is a very long commit message that exceeds...",
		},
		{
			name:     "maxLen zero returns empty",
			input:    "anything",
			maxLen:   0,
			expected: "",
		},
		{
			name:     "maxLen negative returns empty",
			input:    "anything",
			maxLen:   -5,
			expected: "",
		},
		{
			name:     "maxLen 1 truncates without ellipsis",
			input:    "hello",
			maxLen:   1,
			expected: "h",
		},
		{
			name:     "maxLen 3 truncates without ellipsis",
			input:    "hello",
			maxLen:   3,
			expected: "hel",
		},
		{
			name:     "unicode string truncated by rune count",
			input:    "日本語テスト",
			maxLen:   5,
			expected: "日本...",
		},
		{
			name:     "unicode string unchanged when fits",
			input:    "日本語",
			maxLen:   5,
			expected: "日本語",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.input, tt.maxLen)
			if got != tt.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.expected)
			}
		})
	}
}

func TestSprintStatus(t *testing.T) {
	sprint := &domain.Sprint{
		Name: "Feature Sprint",
		Tickets: []domain.Ticket{
			{
				Name:   "auth-flow",
				Branch: "feat/auth",
				Tasks: []domain.Task{
					{Description: "Setup auth"},
					{Description: "Add login"},
					{Description: "Add logout"},
				},
			},
			{
				Name:   "dashboard",
				Branch: "feat/dashboard",
				Tasks: []domain.Task{
					{Description: "Create layout"},
					{Description: "Add widgets"},
				},
			},
			{
				Name:   "settings",
				Branch: "feat/settings",
				Tasks: []domain.Task{
					{Description: "User prefs"},
				},
			},
		},
	}

	tests := []struct {
		name        string
		state       *domain.State
		mustHave    []string
		mustNotHave []string
	}{
		{
			name:  "at start",
			state: &domain.State{CurrentTicket: 0, CurrentTask: 0},
			mustHave: []string{
				`Sprint: "Feature Sprint"`,
				"Progress: 0/3 tickets, 0/6 tasks",
				"Current: Ticket 1 (auth-flow) > Task 1/3",
			},
		},
		{
			name:  "middle of first ticket",
			state: &domain.State{CurrentTicket: 0, CurrentTask: 2},
			mustHave: []string{
				"Progress: 0/3 tickets, 2/6 tasks",
				"Current: Ticket 1 (auth-flow) > Task 3/3",
			},
		},
		{
			name:  "second ticket",
			state: &domain.State{CurrentTicket: 1, CurrentTask: 1},
			mustHave: []string{
				"Progress: 1/3 tickets, 4/6 tasks",
				"Current: Ticket 2 (dashboard) > Task 2/2",
			},
		},
		{
			name:  "completed",
			state: &domain.State{CurrentTicket: 3, CurrentTask: 0},
			mustHave: []string{
				"Progress: 3/3 tickets, 6/6 tasks",
				"Current: Complete",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SprintStatus(sprint, tt.state)
			for _, want := range tt.mustHave {
				if !strings.Contains(got, want) {
					t.Errorf("SprintStatus() missing %q\nGot:\n%s", want, got)
				}
			}
		})
	}
}

func TestCalculateProgress(t *testing.T) {
	sprint := &domain.Sprint{
		Tickets: []domain.Ticket{
			{Tasks: []domain.Task{{}, {}, {}}},
			{Tasks: []domain.Task{{}, {}}},
			{Tasks: []domain.Task{{}}},
		},
	}

	tests := []struct {
		name        string
		state       *domain.State
		wantTickets int
		wantTasks   int
		wantTotal   int
	}{
		{
			name:        "at start",
			state:       &domain.State{CurrentTicket: 0, CurrentTask: 0},
			wantTickets: 0,
			wantTasks:   0,
			wantTotal:   6,
		},
		{
			name:        "mid first ticket",
			state:       &domain.State{CurrentTicket: 0, CurrentTask: 2},
			wantTickets: 0,
			wantTasks:   2,
			wantTotal:   6,
		},
		{
			name:        "first ticket done",
			state:       &domain.State{CurrentTicket: 1, CurrentTask: 0},
			wantTickets: 1,
			wantTasks:   3,
			wantTotal:   6,
		},
		{
			name:        "all done",
			state:       &domain.State{CurrentTicket: 3, CurrentTask: 0},
			wantTickets: 3,
			wantTasks:   6,
			wantTotal:   6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tickets, tasks, total := calculateProgress(sprint, tt.state)
			if tickets != tt.wantTickets {
				t.Errorf("tickets = %d, want %d", tickets, tt.wantTickets)
			}
			if tasks != tt.wantTasks {
				t.Errorf("tasks = %d, want %d", tasks, tt.wantTasks)
			}
			if total != tt.wantTotal {
				t.Errorf("total = %d, want %d", total, tt.wantTotal)
			}
		})
	}
}

func TestPrintSprintComplete(t *testing.T) {
	sprint := &domain.Sprint{
		Name: "Test Sprint",
		Tickets: []domain.Ticket{
			{Tasks: []domain.Task{{}, {}}},
			{Tasks: []domain.Task{{}}},
		},
	}
	state := &domain.State{CurrentTicket: 2, CurrentTask: 0}

	t.Run("outputs completion message", func(t *testing.T) {
		config.SetPlain(true)
		defer config.ResetPlain()
		output := testutil.CaptureStdout(t, func() {
			PrintSprintComplete(sprint, state)
		})
		testutil.AssertContains(t, output, "Test Sprint")
		testutil.AssertContains(t, output, "complete")
		testutil.AssertContains(t, output, "2 tickets")
		testutil.AssertContains(t, output, "3 tasks")
	})
}

func TestPrintSprintStuck(t *testing.T) {
	config.SetPlain(true)
	defer config.ResetPlain()

	sprint := &domain.Sprint{
		Name: "Test Sprint",
		Tickets: []domain.Ticket{
			{
				Name: "auth",
				Tasks: []domain.Task{
					{Description: "Setup auth"},
					{Description: "Add login"},
				},
			},
		},
	}

	t.Run("stuck on task", func(t *testing.T) {
		state := &domain.State{
			CurrentTicket: 0,
			CurrentTask:   1,
			FailureCount:  3,
		}
		output := testutil.CaptureStderr(t, func() {
			PrintSprintStuck(sprint, state)
		})
		testutil.AssertContains(t, output, "stuck")
		testutil.AssertContains(t, output, "3 failures")
		testutil.AssertContains(t, output, "Add login")
	})

	t.Run("stuck after completion", func(t *testing.T) {
		state := &domain.State{
			CurrentTicket: 1,
			CurrentTask:   0,
			FailureCount:  3,
		}
		output := testutil.CaptureStderr(t, func() {
			PrintSprintStuck(sprint, state)
		})
		testutil.AssertContains(t, output, "stuck after completion")
	})
}
