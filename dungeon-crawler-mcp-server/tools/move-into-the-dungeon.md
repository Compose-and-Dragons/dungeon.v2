# Move Into The Dungeon Tool

> 🚧 this is a work in progress

## Overview

This module provides tools for player movement within the dungeon game. It handles two main movement mechanisms:
- `move_by_direction`: Basic directional movement
- `move_player`: Primary navigation tool for dungeon exploration

## Tool Definitions

### move_by_direction
- **Purpose**: Move player in specified direction (north, south, east, west)
- **Usage**: Try "move by north"
- **Parameter**: `direction` (required string)

### move_player  
- **Purpose**: Primary navigation tool for exploring rooms
- **Usage**: "move player north" or "go east"
- **Parameter**: `direction` (required string, must be lowercase cardinal direction)

## Movement Flow

```mermaid
flowchart TD
    A[Player requests movement] --> B[Validate player exists]
    B -->|No player| C[Return error: No player exists]
    B -->|Player exists| D[Get direction argument]
    D --> E[Calculate new coordinates]
    E --> F[Validate boundaries]
    F -->|Out of bounds| G[Return boundary error]
    F -->|Valid position| H[Update player position]
    H --> I[Update player room ID]
    I --> J[Find existing room]
    J -->|Room exists| K[Mark room as visited]
    J -->|New room| L[Generate new room]
    L --> M[Add NPCs/Monsters/Items]
    M --> N[Add room to dungeon]
    K --> O[Return success message]
    N --> O
```


