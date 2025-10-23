package agents

import (
	"fmt"
	"npc-agent-with-mcp/helpers"

	"github.com/firebase/genkit/go/ai"
)

func displayToolsList(tools []ai.ToolRef) {
	shouldIDisplay := helpers.GetEnvOrDefault("LOG_TOOL_MESSAGES", "true")
	if helpers.StringToBool(shouldIDisplay) {
		fmt.Println("🛠️ Tools index", len(tools), "active tools.")
		for _, t := range tools {
			fmt.Println("   -", t.Name())
		}
		fmt.Println()
	}
}
