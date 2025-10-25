

```golang
dungeonMasterToolsAgent.DetectAndExecuteToolCalls(ctx, config, `
    Create a Rogue Dwarf Warrior character named Bob.
`)

fmt.Println(strings.Repeat("=", 50))

dungeonMasterToolsAgent.DetectAndExecuteToolCalls(ctx, config, `
    I want to speak to Elara.
`)

sorcererAgent.LoopCompletion(ctx, npcagents.GetSorcererAgentConfig())

stringResult, errResult := dungeonMasterToolsAgent.DirectExecuteTool(ctx, config, &ai.ToolRequest{
    Name: "c&d_is_player_in_same_room_as_npc",
    Input: map[string]any{
        "name": "Elara",
    },
    Ref: "xxx",
})
if errResult != nil {
    log.Fatal("😡:", errResult)
}
fmt.Println("🟢🟢🟢 🛠️ DirectExecuteTool Result:", stringResult)
fmt.Println(strings.Repeat("=", 50))
```