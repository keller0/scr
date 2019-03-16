#!/bin/sh

export YXI_BACK_PORT=":8090"
export GIN_MODE="debug"
export GIN_LOG_PATH="./gin.log"

export CONTAINER_MEM_LIMIT="50"
export CONTAINER_DISK_LIMIT="5"
export MAX_PUBLIC_RUN=20