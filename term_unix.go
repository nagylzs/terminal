//go:build !windows
// +build !windows

package terminal

import (
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

func (t *Terminal) updatePTYSize() {
	if t.pty == nil { // SSH or other direct connection?
		return
	}
	_ = pty.Setsize(t.pty.(*os.File), &pty.Winsize{
		Rows: uint16(t.config.Rows), Cols: uint16(t.config.Columns),
		X: uint16(t.config.Width), Y: uint16(t.config.Height)})
}

func (t *Terminal) startPTY() (io.WriteCloser, io.Reader, io.Closer, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash"
	}

	_ = os.Chdir(t.startingDir())
	env := os.Environ()
	env = append(env, "TERM=xterm-256color")
	c := exec.Command(shell)
	c.Env = env
	t.cmd = c

	// Start the command with a pty.
	f, err := pty.Start(c)
	return f, f, f, err
}
