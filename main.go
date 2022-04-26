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
	cfg["notification.smtp"] = json.Get("notification.smtp").String()
	cfg["notification.port"] = json.Get("notification.port").String()
	cfg["notification.from"] = json.Get("notification.from").String()
	cfg["notification.to"] = json.Get("notification.to").String()
	cfg["notification.username"] = json.Get("notification.username").String()
	cfg["notification.password"] = json.Get("notification.password").String()

	return cfg, nil
}

func main() {
	fmt.Println("[hushbackup] Start running.")

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "No valid configuration file provided!")
		os.Exit(1)
	}

	// Load config.
	cfgPath := string(os.Args[1])
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config file! Got err: %v\n", err)
		panic(err)
	}

	// Tarball local folder.
	tarFullpath, err := TarFolder(cfg["source.tarballFolder"], cfg["source.tempFolder"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to tarball the folder! The err is: %v\n", err)
		panic(err)
	}

	// Connect to target.
	port, _ := strconv.Atoi(cfg["target.port"])
	// sc means SFTP Connection
	sftpc, sshc, err := ConnectSFTPServer(cfg["target.host"], cfg["target.username"], cfg["target.password"], port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to host! Got err: %v\n", err)
		panic(err)
	}
	defer sshc.Close()
	defer sftpc.Close()

	// Upload file to remote.
	tarRemotePath := filepath.Join(cfg["target.saveFolder"], filepath.Base(tarFullpath))
	err = UploadFile(sftpc, tarFullpath, tarRemotePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to upload to host! Got err: %v\n", err)
		panic(err)
	}

	if len(cfg["notification.to"]) > 0 {
		port, _ = strconv.Atoi(cfg["notification.port"])
		err = SendEMailNotification(cfg["notification.from"], cfg["notification.to"], tarRemotePath, cfg["notification.smtp"], port, cfg["notification.username"], cfg["notification.password"])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send notification email! Got err: %v\n", err)
			// but we dont panic
		}
		fmt.Println("Notification email sent.")
	}

	fmt.Println("Run all done successfully~")
}
