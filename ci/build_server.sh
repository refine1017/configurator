#!/usr/bin/env bash

ls -l
ls -l ./web/dist
yes | cp ./web/dist/bindata.go ./server/public/bindata.go
sed -i 's/main/public/' ./server/public/bindata.go
mkdir package
mkdir package/bin
cd server
go mod download
go build -o ../package/bin/server
cd ..
ls -l