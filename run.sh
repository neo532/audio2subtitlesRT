#/bin/sh
go build main.go && ffmpeg -f avfoundation -i :0 -f segment -segment_time 5 -ar 16000 -f s16le -ac 1 pipe:1|./main
