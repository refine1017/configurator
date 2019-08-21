#!/bin/bash

pwd
ls -l
mv package/bin/server /data/configurator/bin/server
sudo supervisorctl restart configurator
