#!/bin/zsh
echo "[INFO]: Building ./bin/SSHoneyNet and outputting to ./cmd"
go build -o ./bin/SSHoneyNet ./cmd
echo "[INFO]: Running..."
./bin/SSHoneyNet