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

// Connect to remote SFTP server.
func ConnectSFTPServer(host string, user string, password string, port int) (*sftp.Client, *ssh.Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Fprintf(os.Stdout, "Connecting to %s ...\n", addr)

	auths := []ssh.AuthMethod{ssh.Password(password)}
	config := ssh.ClientConfig{
		User:            user,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshc, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return nil, nil, err
	}

	sftpc, err := sftp.NewClient(sshc)
	if err != nil {
		return nil, nil, err
	}

	return sftpc, sshc, nil
}

// Upload certain file via SFTP connection.
func UploadFile(sc *sftp.Client, localFile, remoteFile string) error {
	fmt.Fprintf(os.Stdout, "Upload [%s] to [%s]\n", localFile, remoteFile)

	srcFile, err := os.Open(localFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open localFile: %v\n", err)
		return err
	}
	defer srcFile.Close()

	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		sc.Mkdir(path)
	}

	dstFile, err := sc.OpenFile(remoteFile, (os.O_WRONLY | os.O_CREATE | os.O_TRUNC))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open remoteFile: %v\n", err)
		return err
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to upload: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%.2f MBytes transferred.\n", float32(bytes)/(1024*1024))

	return nil
}
