```bash
# at the root of the repo
#docker compose up --build -d
docker compose up -d
#docker compose logs mcp-gateway

docker compose attach dungeon-master
# then start: ./dungeon-master
# or
docker compose exec -it dungeon-master ./dungeon-master
```