package execution

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	Language string
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
		var cmd *exec.Cmd
		if payload.Branch != "" {
			cmd = exec.CommandContext(ctx, "git", "clone", "--branch", payload.Branch, "--single-branch", payload.RepoURL, ".")
		} else {
			cmd = exec.CommandContext(ctx, "git", "clone", "--single-branch", payload.RepoURL, ".")
		}
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

	// 2.5 Auto-Detect Project Type if commands are empty
	if len(payload.Dependencies) == 0 && len(payload.BuildCmds) == 0 && len(payload.TestCmds) == 0 {
		var packageJsonPath, goModPath, rustTomlPath, pomXmlPath, pyReqPath string
		var hasGo, hasPython bool

		filepath.Walk(workspace, func(path string, info os.FileInfo, err error) error {
			if err != nil { return nil }
			if !info.IsDir() {
				name := info.Name()
				if name == "package.json" && packageJsonPath == "" { packageJsonPath = path }
				if name == "go.mod" && goModPath == "" { goModPath = path }
				if name == "Cargo.toml" && rustTomlPath == "" { rustTomlPath = path }
				if name == "pom.xml" && pomXmlPath == "" { pomXmlPath = path }
				if name == "requirements.txt" && pyReqPath == "" { pyReqPath = path }
				if filepath.Ext(name) == ".go" { hasGo = true }
				if filepath.Ext(name) == ".py" { hasPython = true }
			}
			return nil
		})

		getPrefix := func(fullPath string) string {
			dir := filepath.Dir(fullPath)
			rel, _ := filepath.Rel(workspace, dir)
			if rel != "." && rel != "" {
				return fmt.Sprintf("cd %s && ", rel)
			}
			return ""
		}

		if packageJsonPath != "" {
			payload.Language = "Node.js"
			prefix := getPrefix(packageJsonPath)
			payload.Dependencies = []string{prefix + "npm install"}
			payload.BuildCmds = []string{prefix + "npm run build --if-present"}
			payload.TestCmds = []string{prefix + "npm run test --if-present"}
		} else if goModPath != "" {
			payload.Language = "Go"
			prefix := getPrefix(goModPath)
			payload.Dependencies = []string{prefix + "go mod download"}
			payload.BuildCmds = []string{prefix + "go build ./..."}
			payload.TestCmds = []string{prefix + "go test ./..."}
		} else if rustTomlPath != "" {
			payload.Language = "Rust"
			prefix := getPrefix(rustTomlPath)
			payload.BuildCmds = []string{prefix + "cargo build"}
			payload.TestCmds = []string{prefix + "cargo test"}
		} else if pomXmlPath != "" {
			payload.Language = "Java"
			prefix := getPrefix(pomXmlPath)
			payload.BuildCmds = []string{prefix + "mvn compile"}
			payload.TestCmds = []string{prefix + "mvn test"}
		} else if pyReqPath != "" {
			payload.Language = "Python"
			prefix := getPrefix(pyReqPath)
			payload.Dependencies = []string{prefix + "pip install -r requirements.txt"}
			payload.TestCmds = []string{prefix + "pytest"}
		} else if hasGo {
			payload.Language = "Go (Simple)"
			payload.BuildCmds = []string{"go build"}
			payload.TestCmds = []string{"go test"}
		} else if hasPython {
			payload.Language = "Python (Simple)"
			payload.TestCmds = []string{"python3 -m unittest discover"}
		}
	}
	res.Language = payload.Language
	
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

	if len(payload.Dependencies) == 0 && len(payload.BuildCmds) == 0 && len(payload.TestCmds) == 0 {
		res.Error = fmt.Errorf("no commands executed: could not auto-detect language and no commands were configured")
		res.ExitCode = 1
		res.Duration = time.Since(start)
		return res
	}

	res.Duration = time.Since(start)
	res.ExitCode = 0
	return res
}
