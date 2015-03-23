go clean
GOARM=7 GOARCH=arm GOOS=linux go build
scp horus-v2 root@192.168.7.2:/usr/bin/
