package print

import (
  "os/exec"
  "bufio"
  "fmt"
  "io"
)

type streamer struct {
  prefix string
}

func (s streamer) Write(p []byte) (n int, err error) {
  fmt.Printf("%s%s", s.prefix, p)
  return len(p), nil
}

// Stream executes a pre-assembled command and streams the output with a prefix
func NewStreamer(prefix string) io.Writer {
  return streamer{prefix}
}

func Stream(cmd *exec.Cmd, prefix string) error {

  // setup stderr pipe
  stderr, err := cmd.StderrPipe()
  if err != nil {
    return err
  }

  // create a stderr scanner
  stderrScanner := bufio.NewScanner(stderr)
	go func() {
    // scan lines and print them with a prefix
		for stderrScanner.Scan() {
			fmt.Printf("%s%s\n", prefix, stderrScanner.Text())
		}
	}()

  // setup stdout pipe
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    return err
  }

  // create a stdout scanner
  stdoutScanner := bufio.NewScanner(stdout)
  go func() {
    // scan lines and print them with a prefix
    for stdoutScanner.Scan() {
      fmt.Printf("%s%s\n", prefix, stdoutScanner.Text())
    }
  }()

  // start the command
  if err := cmd.Start(); err != nil {
    return err
  }

  // wait for command to finish
  return cmd.Wait()
}