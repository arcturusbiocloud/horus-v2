go clean
GOARM=7 GOARCH=arm GOOS=linux go build
scp horus-v2 root@10.1.10.111:/usr/bin/
