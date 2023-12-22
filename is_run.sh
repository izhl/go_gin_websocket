#!/bin/bash
# 探针，检测进程是否存在
count=`ps aux | grep go_gin_websocket | wc -l`

if [[ $count>1 ]]; then
#    echo 0
    exit 0
else
    # echo 1
    exit 1
fi