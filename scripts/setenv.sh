#!/bin/sh

export YXI_BACK_PORT=":8090"
export GIN_MODE="debug"
export GIN_LOG_PATH="./gin.log"
#JWT KEY
export YXI_BACK_KEY="secretkey"
export YXI_BACK_MYSQL_ADDR="127.0.0.1:3306"
export YXI_BACK_MYSQL_NAME="yxi"
export YXI_BACK_MYSQL_USER="root"
export YXI_BACK_MYSQL_PASS="111"