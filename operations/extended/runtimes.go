package extended

import (
	"context"

	"github.com/dtyg123/mcp-1panel-extended/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	ListRuntimes = "list_runtimes"
)

var ListRuntimesTool = mcp.NewTool(
	ListRuntimes,
	mcp.WithDescription("list all runtimes (php, go, node, python, java, dotnet) with version, type, status and container image. Note: runtime versions are tracked by 1Panel and may differ from the actual container image tag."),
	mcp.WithNumber("page", mcp.Description("page number, default 1"), mcp.DefaultNumber(1)),
	mcp.WithNumber("pageSize", mcp.Description("page size, default 100"), mcp.DefaultNumber(100)),
)

func ListRuntimesHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	req := map[string]interface{}{
		"page":     int(numArg(args, "page", 1)),
		"pageSize": int(numArg(args, "pageSize", 100)),
	}
	res := map[string]interface{}{}
	return utils.NewPanelClient("POST", "/runtimes/search", utils.WithPayload(req)).Request(&res)
}
