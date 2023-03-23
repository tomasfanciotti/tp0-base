

python3 tp0-doccgen.py --clients $1
docker compose -f docker-compose-gen.yaml up -d --build
