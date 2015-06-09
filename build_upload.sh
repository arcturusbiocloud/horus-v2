go clean
GOARM=7 GOARCH=arm GOOS=linux go build

# copy additional files
#scp -r ~/Documents/arcturusbiocloud/horus-v2/bin root@192.168.7.42:/root/horus-v2/
#scp -r ~/Documents/arcturusbiocloud/horus-v2/streaming root@192.168.7.42:/root/horus-v2/streaming

# copy just the executable
scp ~/Documents/arcturusbiocloud/horus-v2/horus-v2 root@192.168.7.42:/root/horus-v2/
