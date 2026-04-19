package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rusintez/pm/internal/config"
)

func Run(cfg *config.Config, session string, names []string, layoutEnv map[string]string) error {
	inTmux := os.Getenv("TMUX") != ""

	for i, name := range names {
		proc, ok := cfg.Processes[name]
		if !ok {
			return fmt.Errorf("unknown process: %q", name)
		}
		fullCmd := buildCmd(proc.Cmd, layoutEnv, proc.Env)

		if i == 0 {
			if inTmux {
				if err := run("tmux", "new-window", "-n", name, fullCmd); err != nil {
					return err
				}
			} else {
				if err := run("tmux", "new-session", "-d", "-s", session, "-n", name, fullCmd); err != nil {
					return err
				}
			}
		} else {
			if err := run("tmux", "new-window", "-t", session, "-n", name, fullCmd); err != nil {
				return err
			}
		}
	}

	if !inTmux {
		return run("tmux", "attach-session", "-t", session)
	}
	return nil
}

func Attach(session string) error {
	return run("tmux", "attach-session", "-t", session)
}

func Stop(session string) error {
	return run("tmux", "kill-session", "-t", session)
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildCmd(cmdStr string, layoutEnv, procEnv map[string]string) string {
	merged := make(map[string]string)
	for k, v := range layoutEnv {
		merged[k] = v
	}
	for k, v := range procEnv {
		merged[k] = v
	}
	if len(merged) == 0 {
		return cmdStr
	}
	var parts []string
	for k, v := range merged {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, " ") + " " + cmdStr
}
