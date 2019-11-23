package utils

// Cmd abstracts over running a command somewhere, this is useful for testing
type Cmd interface {
	// Run executes the command (like os/exec.Cmd.Run), it should return
	Run() error
	// Each entry should be of the form "key=value"
	Exec() error
}
