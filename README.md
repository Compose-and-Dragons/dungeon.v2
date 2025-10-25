# Dungeon.v2
> With Genkit, Docker Model Runner, Docker Agentic Compose and Docker MCP Gateway.


```bash
docker compose up --build -d
docker compose logs mcp-gateway

docker compose attach dungeon-master
# then start: ./dungeon-master
# or
docker compose exec -it dungeon-master ./dungeon-master

```