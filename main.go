//
// hushbackup is a tool that backup your folder into another server via SFTP
//
// by z0gSh1u @ github.com/z0gSh1u/hushbackup
//
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/tidwall/gjson"
)

func loadConfig(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	json := gjson.Parse(string(data))

	cfg := make(map[string]string)
	cfg["source.tarballFolder"] = json.Get("source.tarballFolder").String()
	cfg["source.tempFolder"] = json.Get("source.tempFolder").String()
	cfg["target.host"] = json.Get("target.host").String()
	cfg["target.port"] = json.Get("target.port").String()
	cfg["target.username"] = json.Get("target.username").String()
	cfg["target.method"] = json.Get("target.method").String()
	cfg["target.password"] = json.Get("target.password").String()
	cfg["target.saveFolder"] = json.Get("target.saveFolder").String()

	return cfg, nil
}

func main() {
	// Load config.
	cfgPath := string(os.Args[1])
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config file! Got err: %v\n", err)
		os.Exit(1)
	}

	// Tarball local folder.
	tarFullpath, err := TarFolder(cfg["source.tarballFolder"], cfg["source.tempFolder"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to tarball the folder! The err is: %v\n", err)
		os.Exit(1)
	}

	// Connect to target.
	port, err := strconv.Atoi(cfg["target.port"])
	// sc means SFTP Connection
	sc, err := ConnectSFTPServer(cfg["target.host"], cfg["target.username"], cfg["target.password"], port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to host! Got err: %v\n", err)
		os.Exit(1)
	}
	defer sc.Close()

	// Upload file to remote.
	err = UploadFile(sc, tarFullpath, filepath.Join(cfg["target.saveFolder"], filepath.Base(tarFullpath)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to upload to host! Got err: %v\n", err)
		os.Exit(1)
	}

	// TODO

}
