

```bash
docker compose up --build -d

docker compose attach sorcerer-json-agent
# then start: ./sorcerer-agent
# or
docker compose exec -it sorcerer-json-agent ./sorcerer-agent

```