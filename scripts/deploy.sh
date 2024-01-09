#!/usr/bin/env bash
set -x

REMOTE_DIR="/home/xrdpuser/Desktop/vidre-back"
REMOTE_HOST="root@185.166.212.43"

# SSH and execute commands on remote server
if ssh $REMOTE_HOST "cd $REMOTE_DIR && git pull"; then
    echo "Deployment successful."
else
    echo "Deployment failed."
    exit 1
fi