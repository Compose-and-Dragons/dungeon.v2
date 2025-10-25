# Dungeon MCP Server

## Available Tools

### MCP Server Tools (dungeon-crawler-mcp-server)

- `create_player`: Create a new player with name, class, and race. Try: "I'm Bob, the Dwarf Warrior."
- `get_player_info`: Get the current player's information. Try: "Who am I?"
- `get_dungeon_info`: Get the current dungeon's information including its layout, rooms, entrance and exit coordinates
- `get_current_room_info`: Get information about the current room where the player is located. Try: "Where am I?" or "Look around"
- `move_by_direction`: Move the player in a specified direction (north, south, east, west). Try "move by north"
- `move_player`: Move the player in the dungeon by specifying a cardinal direction. This is the primary navigation tool for exploring rooms. Usage: "move player north" or "go east"
- `get_dungeon_map`: Generate an ASCII map of the discovered dungeon rooms showing the player position, NPCs, and monsters with a legend