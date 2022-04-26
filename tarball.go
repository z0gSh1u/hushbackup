package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

// Tarball the folder into .tar archive.
func TarFolder(path string, temp string) (fullpath string, err error) {
	timeNow := time.Now().Format("20060102150405") // stupid design
	tarFilename := filepath.Base(path) + "." + timeNow + ".tar"
	tarFullpath := filepath.Join(temp, tarFilename)

	cmd := exec.Command("tar", "cvf", filepath.Join(temp, tarFilename), path)
	fmt.Printf("Tarball command: %s\n", cmd.String())
	_, err = cmd.Output()
	if err != nil {
		return "", err
	}

	return tarFullpath, nil
}
