package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger"
)

//FileModeStandardFolder owner can rwx, else r-x
//https://stackoverflow.com/questions/14249467/os-mkdir-and-os-mkdirall-permission-value
const FileModeStandardFolder = 0755

//CopyFromPkgr copies a single file or a directory with
//all of its contents to the destination directory.
func CopyFromPkgr(pkgrSrc string, dest string) error {
	var destDirCreated = false
	var basePath string
	err := pkger.Walk(pkgrSrc, func(path string, info os.FileInfo, err error) error {
		if !destDirCreated {
			destDirCreated = true
			if info.IsDir() {
				os.MkdirAll(dest, FileModeStandardFolder)
				basePath = path
			} else {
				os.MkdirAll(filepath.Dir(dest), FileModeStandardFolder)
				copyPlainFromPkgr(path, dest)
			}
			return nil
		}

		destPath := fmt.Sprintf("%s%s", dest, strings.TrimPrefix(path, basePath))

		if info.IsDir() {
			os.MkdirAll(destPath, FileModeStandardFolder)
		} else {
			copyPlainFromPkgr(path, destPath)
		}
		return nil
	})
	return err
}

func copyPlainFromPkgr(pkgerSrcPath string, destPath string) error {
	srcFile, err := pkger.Open(pkgerSrcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	return err
}

//CopyPlain copies a plain file.
//too lazy to do full impl
func CopyPlain(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	return err
}

//Unzip unzips files there's a large limit on unzip
func Unzip(pwd string, zipPath string, outputDir string) error {
	cmd := exec.Command(
		"sunzip-cli",
		zipPath,
		"-ms", "15",
		"-mm", "10240",
		"-md", "1024000",
		"-d", outputDir,
	)
	cmd.Env = []string{}
	cmd.Dir = pwd
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
