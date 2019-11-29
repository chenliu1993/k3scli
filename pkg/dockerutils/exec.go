package dockerutils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"bytes"
	"bufio"
	log "github.com/sirupsen/logrus"
)

// Exec uses docker eec to actually exec a process in a container
func (c *ContainerCmd) Exec(in io.Reader, out, stderr io.Writer) error {
	args := []string{
		"exec",
	}
	if c.Detach {
		args = append(args, "-d")
	} else {
		args = append(args, "-it")
	}
	args = append(args, c.ID)
	args = append(args, c.Args...)
	cmd := exec.Command(c.Command, args...)
	fmt.Print(cmd)
	fmt.Print("\n")
	if in == nil {
		cmd.Stdin = os.Stdin
	} else {
		cmd.Stdin = in
	}
	if out == nil {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = out
	}
	if stderr == nil {
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = stderr
	}
	log.Debug(fmt.Sprintf("begin exec process in container: %s", c.ID))
	err := cmd.Run()
	if err != nil {
		log.Debug(err)
		return err
	}
	return nil
}


// ExecOutput save the output to strings
func ExecOutput(cmd exec.Cmd, includeErr bool) (lines []string,err  error) {
	var buf bytes.Buffer
	if includeErr == true {
		cmd.Stderr = &buf
	}
	cmd.Stdout = &buf
	err = cmd.Run()
	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, err
}

// PrintOutput print buffer output to stdout
func PrintOutput(out []string) {
	for _, line := range out {
		fmt.Printf(line+"\n")
	}
}