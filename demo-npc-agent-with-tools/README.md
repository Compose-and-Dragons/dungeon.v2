

```bash
docker compose up --build -d
docker compose logs mcp-gateway

docker compose attach sorcerer-mcp-agent
# then start: ./sorcerer-agent
# or
docker compose exec -it sorcerer-mcp-agent ./sorcerer-agent
```

Try:
```plaintext
Roll 3 dices with 6 faces each.
Then generate a character name for an elf.
Finally, roll 2 dices with 8 faces each.
After that, generate a character name for a dwarf.
```