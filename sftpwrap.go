package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/sftp"
)

func ConnectSFTPServer(host string, user string, password string, port int) (*sftp.Client, error) {
	fmt.Fprintf(os.Stdout, "Connecting to %s ...\n", host)

	auths := []ssh.AuthMethod{ssh.Password(password)}

	config := ssh.ClientConfig{
		User:            user,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connecto to [%s]: %v\n", addr, err)
		return nil, err
	}
	defer conn.Close()

	sc, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start SFTP subsystem: %v\n", err)
		return nil, err
	}

	defer sc.Close()
	return sc, nil
}

func UploadFile(sc *sftp.Client, localFile, remoteFile string) (err error) {
	fmt.Fprintf(os.Stdout, "Uploading [%s] to [%s] ...\n", localFile, remoteFile)

	srcFile, err := os.Open(localFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open local file: %v\n", err)
		return
	}
	defer srcFile.Close()

	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		sc.Mkdir(path)
	}

	// Note: SFTP To Go doesn't support O_RDWR mode
	dstFile, err := sc.OpenFile(remoteFile, (os.O_WRONLY | os.O_CREATE | os.O_TRUNC))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open remote file: %v\n", err)
		return
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to upload local file: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%d bytes copied\n", bytes)

	return
}
