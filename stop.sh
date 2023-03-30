
docker compose -f docker-compose-gen.yaml stop -t 5

if [ "$1" == "-k" ]; then
  docker compose -f docker-compose-gen.yaml down
fi