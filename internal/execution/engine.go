package execution

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type JobPayload struct {
	RepoURL      string   `json:"repo_url"`
	Branch       string   `json:"branch"`
	Language     string   `json:"language"`
	Dependencies []string `json:"dependencies"`
	BuildCmds    []string `json:"build_cmds"`
	TestCmds     []string `json:"test_cmds"`
}

type ExecutionResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
	Error    error
}

type Engine interface {
	ExecuteJob(ctx context.Context, payloadStr string) *ExecutionResult
}

type engine struct{}

func NewEngine() Engine {
	return &engine{}
}

func getShell() (string, string) {
	if os.PathSeparator == '\\' {
		return "cmd", "/c"
	}
	return "sh", "-c"
}

func (e *engine) ExecuteJob(ctx context.Context, payloadStr string) *ExecutionResult {
	start := time.Now()
	res := &ExecutionResult{}

	var payload JobPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		res.Error = fmt.Errorf("failed to parse payload: %w", err)
		res.ExitCode = -1
		return res
	}

	// 1. Prepare Workspace
	workspace, err := os.MkdirTemp("", "forgeflow-job-*")
	if err != nil {
		res.Error = fmt.Errorf("failed to create workspace: %w", err)
		res.ExitCode = -1
		return res
	}
	defer os.RemoveAll(workspace)

	// 2. Clone/Extract
	if payload.RepoURL != "" {
		branch := payload.Branch
		if branch == "" {
			branch = "main" // default
		}
		cmd := exec.CommandContext(ctx, "git", "clone", "--branch", branch, "--single-branch", payload.RepoURL, ".")
		cmd.Dir = workspace
		if out, err := cmd.CombinedOutput(); err != nil {
			res.Error = fmt.Errorf("git clone failed: %v", err)
			res.Stderr = string(out)
			if exitError, ok := err.(*exec.ExitError); ok {
				res.ExitCode = exitError.ExitCode()
			} else {
				res.ExitCode = 1
			}
			res.Duration = time.Since(start)
			return res
		}
	}

	// Helper to run commands
	shell, arg := getShell()
	runCmds := func(cmds []string, phase string) error {
		for _, c := range cmds {
			cmd := exec.CommandContext(ctx, shell, arg, c)
			cmd.Dir = workspace
			out, err := cmd.CombinedOutput()
			res.Stdout += fmt.Sprintf("=== %s: %s ===\n%s\n", phase, c, string(out))
			if err != nil {
				res.Stderr += fmt.Sprintf("=== %s: %s ===\n%s\n", phase, c, string(out))
				if exitError, ok := err.(*exec.ExitError); ok {
					res.ExitCode = exitError.ExitCode()
				} else {
					res.ExitCode = 1
				}
				return fmt.Errorf("command failed during %s: %s, err: %v", phase, c, err)
			}
		}
		return nil
	}

	// 3. Install dependencies
	if err := runCmds(payload.Dependencies, "Dependencies"); err != nil {
		res.Error = err
		res.Duration = time.Since(start)
		return res
	}

	// 4. Run build
	if err := runCmds(payload.BuildCmds, "Build"); err != nil {
		res.Error = err
		res.Duration = time.Since(start)
		return res
	}

	// 5. Run tests
	if err := runCmds(payload.TestCmds, "Test"); err != nil {
		res.Error = err
		res.Duration = time.Since(start)
		return res
	}

	res.Duration = time.Since(start)
	res.ExitCode = 0
	return res
}
