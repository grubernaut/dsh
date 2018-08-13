package dsh

import (
	"os"
	"sync"
)

//ExecOpts does things
type ExecOpts struct {
	ShowNames         bool
	ShowAddresses     bool
	ShowUsername      bool
	RemoteShell       string
	RemoteUser        string
	RemoteCommand     string
	RemoteCommandOpts string
	ConcurrentShell   bool
	Verbose           bool
	OutputFile        string
}

// Endpoint represents an individual node to connect to.
// Passed to Execute as a slice
type Endpoint struct {
	DisplayName string
	Machine     string
}

// Signal is returned from a goroutine via a channel
type signal struct {
	err       error
	errOutput string
}

type shell struct {
	RemoteCmd     string
	RemoteUser    string
	CmdOpts       []string
	C             chan signal
	ShowNames     bool
	ShowAddresses bool
	ShowUsername  bool
	Node          Endpoint
	Output        *os.File
	fLock         sync.Mutex
}
