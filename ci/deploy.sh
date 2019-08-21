#!/bin/bash

# server config
ssh_user=$DEPLOY_SSH_USER
ssh_host=$DEPLOY_SSH_HOST
ssh_port=$DEPLOY_SSH_PORT
datetime=`date +%Y%m%d%H%M%S`

echo $DEPLOY_SSH_RSA > id_rsa
chmod 600 id_rsa

## stop server
echo "stop server"
expect <<!
set timeout 600
spawn ssh -i id_rsa ${ssh_user}@${ssh_host} -p ${ssh_port}
expect {
 "yes/no" { send "yes\r"; exp_continue}
 "passphrase" { send "\r" }
}
expect "$"
send "sudo supervisorctl stop all\r"
expect "$"
send "exit\r"
expect eof
!

## copy server
echo "copy server"
expect <<!
set timeout 600
spawn scp -i id_rsa -r -P ${ssh_port} ./package/bin/server ${ssh_user}@${ssh_host}:/home/xzhang/configurator_server
expect {
 "yes/no" { send "yes\r"; exp_continue}
 "passphrase" { send "\r" }
}
expect eof
!

## start server
echo "start server"
expect <<!
set timeout 600
spawn ssh -i id_rsa ${ssh_user}@${ssh_host} -p ${ssh_port}
expect {
 "yes/no" { send "yes\r"; exp_continue}
 "passphrase" { send "\r" }
}
expect "$"
send "sudo mv /home/xzhang/configurator_server /data/configurator/bin/server\r"
expect "$"
send "sudo chmod +x /data/configurator/bin/server\r"
expect "$"
send "sudo supervisorctl restart all\r"
expect "$"
send "sudo supervisorctl status\r"
expect "$"
send "exit\r"
expect eof
!
