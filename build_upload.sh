go clean
GOARM=7 GOARCH=arm GOOS=linux go build
#scp -r ~/Documents/arcturusbiocloud/horus-v2/bin root@10.1.10.111:/root/horus-v2/
#scp -r ~/Documents/arcturusbiocloud/horus-v2/streaming root@10.1.10.111:/root/horus-v2/streaming
scp ~/Documents/arcturusbiocloud/horus-v2/horus-v2 root@192.168.7.42:/root/horus-v2/
