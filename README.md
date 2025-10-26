# Dungeon.v2
> With Genkit, Docker Model Runner, Docker Agentic Compose and Docker MCP Gateway.

## Architecture

The project uses an architecture based on:
- **NPC Agents**: Conversational agents with personality and vector memory (RAG)
- **Dungeon Master**: Main agent with tool detection and execution via MCP (Model Context Protocol)
- **Conversation History**: Reset after each interaction to avoid tool call accumulation

### Important Changes

The Dungeon Master **resets its message history** after each interaction (`ResetMessages()` - line 255 in [dungeon-master/main.go](dungeon-master/main.go#L255)). This prevents:
- Accumulation of tool calls in the history
- Unwanted cross-references between interactions
- Overly long contexts that degrade performance

Other agents (Guard, Sorcerer, Healer, Merchant, Boss) maintain their history to preserve conversational coherence with the player.

## Getting Started

```bash
docker compose up --build -d
docker compose logs mcp-gateway

docker compose attach dungeon-master
# then start: ./dungeon-master
# or
docker compose exec -it dungeon-master ./dungeon-master
```