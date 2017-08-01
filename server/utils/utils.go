package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type ScanFunc func(*bufio.Scanner);

func ExitOnDev() {
	if IsDev() {
		os.Exit(1)
	}
}

func BasicScanFunc(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
}

func ReadBody(r *http.Response, t interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(t)
}

func GetBody(r *http.Response) string {
	defer r.Body.Close()
	contents, _ := ioutil.ReadAll(r.Body)
	return fmt.Sprintf("%s", contents)
}

func FatalRunShell(command string, args []string) {
	err := RunShell(command, args)
	if err != nil {
		os.Exit(1)
	}
}

func RunShell(command string, args []string) error {
	return RunShellWithScan(command, args, BasicScanFunc, BasicScanFunc)
}

func RunShellWithScan(command string, args []string, stdio ScanFunc, stderr ScanFunc) error {
	cmdName := command
	cmd := exec.Command(cmdName, args...)
	fmt.Printf("cmd: %s %s\n", cmdName, strings.Join(args, " "))

	if  stdio != nil {
		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintln(os.Stdout, "Error creating StdoutPipe for Cmd", err)
			return err
		}
		go stdio(bufio.NewScanner(cmdReader))
	}

	if stderr != nil {
		errReader, err := cmd.StderrPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating StderrPipe for Cmd", err)
			return err
		}
		go stderr(bufio.NewScanner(errReader))
	}

	err := cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		return err
	}
	return nil
}
