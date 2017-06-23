package utils

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"os/exec"
	"os"
	"bufio"
)

func ExitOnDev(){
	if IsDev() {
		os.Exit(1)
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
	cmdName := command
	cmd := exec.Command(cmdName, args...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s", scanner.Text())
		}
	}()

	err = cmd.Start()
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
