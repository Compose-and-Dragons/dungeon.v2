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
# 🤖📡 Simple Agent with Docker MCP Gateway 
# & 🐙 Agentic Compose
<div class="mermaid">
graph TD
    A[sorcerer-mcp-agent<br/>Go Application] -->|HTTP :9011| B[mcp-gateway<br/>Docker MCP Gateway]
    
    B -->|HTTP :6060/mcp| C[mcp-fancy-names<br/>MCP Server]
    B -->|HTTP :6060/mcp| D[mcp-hello-world<br/>MCP Server]
    B -->|HTTP :6060/mcp| E[mcp-dnd<br/>MCP Server]
    
    B -.->|catalog.yaml| F[Config<br/>Server Registry]
    
    A -.->|models| G[AI Models<br/>chat/embedding/tools]
    
    
    style A fill:#e1f5ff
    style B fill:#ffe1e1
    style C fill:#e8f5e8
    style D fill:#e8f5e8
    style E fill:#e8f5e8
    style F fill:#fff4e1
    style G fill:#f3e5f5
</div>
