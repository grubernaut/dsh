package dsh

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Execute remote commands for each host
func (e *ExecOpts) Execute(nodes []Endpoint) error {
	signals := make(chan signal)
	for _, node := range nodes {
		cmdOpts := e.buildCmdOpts(node.Machine)
		if e.Verbose {
			fmt.Printf("Dumping parameters passed to exec\n")
			fmt.Printf("%#v\n", cmdOpts)
		}

		// Spawn a new go routine for each node
		node := shell{
			RemoteCmd:     e.RemoteShell,
			RemoteUser:    e.RemoteUser,
			CmdOpts:       cmdOpts,
			C:             signals,
			ShowNames:     e.ShowNames,
			ShowAddresses: e.ShowAddresses,
			ShowUsername:  e.ShowUsername,
			Node:          node,
		}

		go node.executeShell()
	}

	// Block until routines are cleaned up
	var err error
	for i := 0; i < len(nodes); i++ {
		select {
		case signal := <-signals:
			if signal.err != nil {
				fmt.Printf(signal.errOutput)
				err = signal.err
			}
		}
	}
	return err
}

// Build up command options for each node
func (e *ExecOpts) buildCmdOpts(machine string) []string {
	var opts []string
	if e.RemoteCommandOpts != "" {
		// Split remote command opts based on <space> to successfully be sent to cmd.Exec()
		remoteOpts := strings.Split(e.RemoteCommandOpts, " ")
		opts = append(opts, remoteOpts...)
	}

	// TODO: use pointers for these to check for nil value
	if e.RemoteUser != "" {
		opts = append(opts, "-l")
		opts = append(opts, e.RemoteUser)
	}
	opts = append(opts, machine)

	opts = append(opts, e.RemoteCommand)
	return opts
}

// Performs the actual execution of the Remote Shell command
func (s *shell) executeShell() {
	// hopefully you don't need it
	var errOutput bytes.Buffer
	run := exec.Command(s.RemoteCmd, s.CmdOpts...)
	run.Stderr = io.Writer(&errOutput)
	stdout, err := run.StdoutPipe()
	if err != nil {
		s.C <- signal{
			err: err,
		}
		return
	}
	run.Env = os.Environ()

	// Get output prefix
	outputPrefix := ""
	if s.ShowNames {
		outputPrefix = s.Node.DisplayName
		if s.ShowAddresses {
			outputPrefix = fmt.Sprintf("%s(%s)", outputPrefix, s.Node.Machine)
		}
		if s.ShowUsername {
			outputPrefix = fmt.Sprintf("%s@%s", s.RemoteUser, outputPrefix)
		}
	}

	// Create a new scanner from stdout pipe
	scanner := bufio.NewScanner(stdout)
	// While we have stdout to print, print it.
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s%s\n", outputPrefix, scanner.Text())
		}
	}()

	if err := run.Start(); err != nil {
		s.C <- signal{
			err:       err,
			errOutput: fmt.Sprintf("%s%s", outputPrefix, errOutput.String()),
		}
		return
	}

	// Block for command to finish
	if err := run.Wait(); err != nil {
		s.C <- signal{
			err:       err,
			errOutput: fmt.Sprintf("%s%s", outputPrefix, errOutput.String()),
		}
		return
	}

	// Non-error case
	s.C <- signal{nil, ""}
}
