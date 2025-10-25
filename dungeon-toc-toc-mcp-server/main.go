package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	// Create MCP server
	s := server.NewMCPServer(
		"mcp-toc-toc",
		"0.0.0",
	)

	speakToAnAgentTool := mcp.NewTool("speak_to_somebody",
		mcp.WithDescription("Speak to somebody by name"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the person to speak to"),
		),
	)
	s.AddTool(speakToAnAgentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

		type Response struct {
			Name string `json:"name"`
		}

		args := request.GetArguments()

		nameArg, exists := args["name"]
		if !exists || nameArg == nil {
			return mcp.NewToolResultText(`{"error": "missing required parameter 'name'"}`), nil
		}

		name, ok := nameArg.(string)
		if !ok {
			return mcp.NewToolResultText(`{"error": "parameter 'name' must be a string"}`), nil
		}

		response := Response{Name: name}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			return mcp.NewToolResultText(`{"error": "failed to marshal response"}`), nil
		}

		return mcp.NewToolResultText(string(responseJSON)), nil
	})

	// Start the HTTP server
	httpPort := os.Getenv("MCP_HTTP_PORT")
	if httpPort == "" {
		httpPort = "9090"
	}

	log.Println("MCP 👋 Toc Toc 🌍 Server is running on port", httpPort)

	// Create a custom mux to handle both MCP and health endpoints
	mux := http.NewServeMux()

	// Add healthcheck endpoint
	mux.HandleFunc("/health", healthCheckHandler)

	// Add MCP endpoint
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	)

	// Register MCP handler with the mux
	mux.Handle("/mcp", httpServer)

	// Start the HTTP server with custom mux
	log.Fatal(http.ListenAndServe(":"+httpPort, mux))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status": "healthy",
		"server": "mcp-toc-toc-server",
	}
	json.NewEncoder(w).Encode(response)
}
