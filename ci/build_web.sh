#!/usr/bin/env bash

ls -l
export APP_ROOT=`pwd`
export PATH=$PATH:$APP_ROOT/ci
chmod +x ./ci/*
cd web
npm install
npm run build:prod
cd dist
go-bindata-assetfs static/... index.html
cd ../../
ls -l