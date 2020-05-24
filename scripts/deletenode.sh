NAME=${1?Error: no name given}
docker rmi `(docker images */$NAME | head -2 | tail -1 | awk '{print $3}')` -f