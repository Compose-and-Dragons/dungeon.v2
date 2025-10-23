

```bash
docker compose up --build -d
docker compose logs mcp-gateway

docker compose attach sorcerer-mcp-agent
# then start: ./sorcerer-agent
# or
docker compose exec -it sorcerer-mcp-agent ./sorcerer-agent

```