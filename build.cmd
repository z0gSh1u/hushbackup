SET HUSHBACKUPVERSION=1.0

@REM Target: Linux amd64
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o bin/hushbackup-%HUSHBACKUPVERSION%-linux-amd64

@REM Target: Linux i386
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=386
go build -o bin/hushbackup-%HUSHBACKUPVERSION%-linux-i386