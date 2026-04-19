package runner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/fatih/color"
	"github.com/rusintez/pm/internal/config"
)

var palette = []*color.Color{
	color.New(color.FgCyan),
	color.New(color.FgGreen),
	color.New(color.FgYellow),
	color.New(color.FgMagenta),
	color.New(color.FgBlue),
	color.New(color.FgRed),
}

func Run(cfg *config.Config, names []string, layoutEnv map[string]string) error {
	maxLen := 0
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	cmds := make([]*exec.Cmd, 0, len(names))
	var wg sync.WaitGroup

	for i, name := range names {
		proc, ok := cfg.Processes[name]
		if !ok {
			return fmt.Errorf("unknown process: %q", name)
		}

		col := palette[i%len(palette)]
		prefix := col.Sprintf("%-*s |", maxLen, name)

		parts := strings.Fields(proc.Cmd)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Env = buildEnv(layoutEnv, proc.Env)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("start %s: %w", name, err)
		}
		cmds = append(cmds, cmd)

		wg.Add(1)
		go func(r io.Reader) {
			defer wg.Done()
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				fmt.Printf("%s %s\n", prefix, scanner.Text())
			}
		}(io.MultiReader(stdout, stderr))
	}

	// wait for Ctrl-C or all processes to exit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-sig:
		for _, cmd := range cmds {
			if cmd.Process != nil {
				cmd.Process.Signal(syscall.SIGTERM)
			}
		}
		wg.Wait()
	case <-done:
	}

	return nil
}

func buildEnv(layoutEnv, procEnv map[string]string) []string {
	env := os.Environ()
	merged := make(map[string]string)
	for k, v := range layoutEnv {
		merged[k] = v
	}
	for k, v := range procEnv {
		merged[k] = v // process wins
	}
	for k, v := range merged {
		env = append(env, k+"="+v)
	}
	return env
}
