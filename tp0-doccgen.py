from copy import copy, deepcopy
from sys import argv
import yaml

MAX_CLIENTS = 10
DEFAULT_CLIENTS = 2

version = '3.9'

# Templates

server = {
    "container_name": "server",
    "image": "server:latest",
    "entrypoint": "python3 /app/main.py",
    "environment": ["PYTHONUNBUFFERED=1","LOGGING_LEVEL=DEBUG"],
    "networks": ["testing_net"],
    "volumes": ["./server/:/app/"]
}

network = {
    "ipam": {
        "driver": "default",
        "config": [
            {"subnet": "172.25.125.0/24"}
        ]
    }
}

client = {
    "container_name": "client",
    "image": "client:latest",
    "entrypoint": "/app/client",
    "environment": ["LOGGING_LEVEL=DEBUG"],
    "networks": ["testing_net"],
    "depends_on": ["server"],
    "volumes": ["./client/:/app/config/"]
}


def generate(clients):
    config = {}
    services = {"server": server}

    clients = min(clients, MAX_CLIENTS)
    for i in range(clients):
        service_name = "client-" + str(i + 1)
        client_aux = deepcopy(client)
        client_aux["container_name"] = service_name
        client_aux["environment"].append(f"CLI_ID={i + 1}")
        services[service_name] = client_aux

    config["services"] = services
    config["version"] = version
    config["networks"] = {"testing_net": network}

    return config


def main():
    print(argv)
    if len(argv) == 3 and argv[1] == "--clients":
        clients = int(argv[2])
    else:
        clients = DEFAULT_CLIENTS

    config = generate(clients)
    print(config)

    with open("docker-compose-gen.yaml", "w") as docc_file:
        yaml.dump(config, docc_file)


if __name__ == "__main__":
    main()