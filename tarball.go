package main

func TarFolder(path string, temp string, gzip bool) string {
	var cmd string
	cmd = "tar -cvf"
	if gzip {
		cmd = cmd + "z"
	}
	return cmd
}
