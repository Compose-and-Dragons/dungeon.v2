# Demo Queen Pédauque Agent
> - https://hub.docker.com/repository/docker/philippecharriere494/queen-pedauque/general
> - https://hub.docker.com/r/philippecharriere494/queen-pedauque/tags
```bash
docker compose up --build -d

docker compose attach queen-pedauque-agent
# then start: ./queen-pedauque-agent
# or
docker compose exec -it queen-pedauque-agent ./queen-pedauque-agent

```