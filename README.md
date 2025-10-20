# mcp-toolgate
A lightweight proxy for stdio-based MCP servers that filters which tools are registered and exposed to clients.

## Motivation

When using MCP servers with AI agents, having too many tools registered can lead to unpredictable behavior —  
agents may become confused or fail to use the intended tool.  

To address this, `mcp-toolgate` provides a simple way to restrict which tools are exposed to the agent.  
By allowing tools via environment variables, it helps create a more controlled and predictable tool environment.


## Example Configuration

You can refer to the [.vscode/mcp.json](.vscode/mcp.json) file for an example setup.  
This configuration demonstrates how to run an MCP server through `mcp-toolgate` to control which tools are exposed to the client.


```jsonc
{
    "servers": {
        "github": {
            "command": "bash",
            "args": [
                "-lc",
                "docker run -i --rm -e GITHUB_PERSONAL_ACCESS_TOKEN ghcr.io/github/github-mcp-server | go run ./cmd/mcp-toolgate"
            ],
            "env": {
                "ALLOWED_TOOL_NAMES": "add_issue_comment",
                "DEBUG": "",
                "GITHUB_PERSONAL_ACCESS_TOKEN": "${env:GITHUB_TOKEN}",
            }
        }
    }
}
```

### Behavior

When used as shown above, `mcp-toolgate` acts as a proxy that intercepts the MCP server’s tool registration and filters it based on environment variables.

### Environment Variables
- `ALLOWED_TOOL_NAMES`: Comma-separated list of tools to allow. Only these tools will be registered and exposed to clients.
- `DEBUG`: Enables debug logging when set (even to an empty string). Useful for troubleshooting.

### Notes

- `mcp-toolgate` reads from stdin and writes to stdout, making it easy to insert into an existing pipeline:
    ```
    mcp-server | mcp-toolgate
    ```
- Intended for simple, configuration-driven filtering of tools without modifying the original MCP server.
