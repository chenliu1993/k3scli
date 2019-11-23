package dockerutils

import (
	"fmt"
	"os"
	"os/exec"
	"bytes"
	"bufio"
	log "github.com/sirupsen/logrus"
)

// Exec uses docker eec to actually exec a process in a container
func (c *ContainerCmd) Exec() error {
	args := []string{
		"exec",
	}
	if c.Detach {
		args = append(args, "-d")
		args = append(args, c.ID)
		args = append(args, c.Args...)
		cmd := exec.Command(c.Command, args...)
		lines, err := ExecOutput(*cmd, true)
		if err != nil {
			return err
		}
		PrintOutput(lines)
	}else{
		args = append(args, "-it")
		args = append(args, c.ID)
		args = append(args, c.Args...)
		cmd := exec.Command(c.Command, args...)
		// applies to the docker -i options
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Debug(fmt.Sprintf("begin exec process in container: %s", c.ID))
		err := cmd.Run()
		if err != nil {
			log.Debug(err)
			return err
		}
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