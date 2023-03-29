
if [ "$1" == "server" ]; then
  docker build -f ./server/Dockerfile -t "server:latest" .

elif [ "$1" == "client" ]; then
  docker build -f ./client/Dockerfile -t "client:latest" .

elif [ "$1" == "" ]; then
    docker build -f ./server/Dockerfile -t "server:latest" .
    docker build -f ./client/Dockerfile -t "client:latest" .
fi
