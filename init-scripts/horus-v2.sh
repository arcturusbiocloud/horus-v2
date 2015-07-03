#! /bin/bash
# /etc/init.d/horus-v2

OUT_LOG=/dev/null
HORUS_BINARY=/home/pi/horus-v2/horus-v2

case "$1" in
  start)
    echo "Starting horus-v2"
    PORT=3001 $HORUS_BINARY 1>$OUT_LOG 2>$OUT_LOG &
    ;;
  stop)
    echo "Stopping horus-v2"
    killall horus-v2
    ;;
  *)
    echo "Usage: /etc/init.d/horus-v2 {start|stop}"
    exit 1
    ;;
  esac
exit 0