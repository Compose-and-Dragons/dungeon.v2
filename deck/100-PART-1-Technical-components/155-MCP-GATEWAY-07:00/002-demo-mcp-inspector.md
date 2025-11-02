---
marp: true
html: true
theme: default
paginate: true
---
<style>
.dodgerblue {
  color: dodgerblue;
}
</style>
# 📡 Docker MCP Gateway & 🕵️‍♂️ Inspector

## How to Create a MCP Server
> <span class="dodgerblue">[/demo-npc-agent-with-tools/mcp-servers/hello-world/main.go](/demo-npc-agent-with-tools/mcp-servers/hello-world/main.go)</span>
## **Demo** (without LLM)
> Using the [MCP Gateway](/demo-npc-agent-with-tools/compose.only.mcp.yaml) and MCP Inspector to connect to and explore MCP servers
```bash
docker compose -f compose.only.mcp.yaml up
```

```bash
# using http://localhost:9011/mcp
npx @modelcontextprotocol/inspector@0.17.2
```