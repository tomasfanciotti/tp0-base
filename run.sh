
if [ "$1" == "test" ]; then

    # Temporary container to test the server
    docker build -f ./netcat/Dockerfile -t "netcat:latest" .
    docker run --rm --network "$(cat ./server/network)" -i -t "netcat:latest"

else

  # Docker Compose Up
  python3 tp0-doccgen.py --clients "$1"
  docker compose -f docker-compose-gen.yaml up -d --build

fi
