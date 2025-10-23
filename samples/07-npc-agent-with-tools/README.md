

```bash
docker compose up --build -d
docker compose attach sorcerer-tools-agent
# then start: ./sorcerer-agent
# or
docker compose exec -it sorcerer-tools-agent ./sorcerer-agent

```