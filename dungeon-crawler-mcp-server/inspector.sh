#!/bin/bash
: <<'COMMENT'
http://localhost:9090/mcp
http://localhost:9011/mcp
COMMENT

npx @modelcontextprotocol/inspector@0.17.2

