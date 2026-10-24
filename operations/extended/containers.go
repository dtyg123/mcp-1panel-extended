package extended

import (
	"context"

	"github.com/dtyg123/mcp-1panel-extended/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	ListContainers = "list_containers"
)

var ListContainersTool = mcp.NewTool(
	ListContainers,
	mcp.WithDescription("list docker containers with name, image, state, status, ports and networks"),
	mcp.WithString("state",
		mcp.Description("filter by container state"),
		mcp.DefaultString("all"),
		mcp.Enum("all", "created", "running", "paused", "restarting", "removing", "exited", "dead")),
	mcp.WithString("name",
		mcp.Description("filter by container name (substring), optional"),
		mcp.DefaultString("")),
	mcp.WithString("orderBy",
		mcp.Description("sort field"),
		mcp.DefaultString("createdAt"),
		mcp.Enum("name", "createdAt", "state")),
	mcp.WithString("order",
		mcp.Description("sort direction"),
		mcp.DefaultString("descending"),
		mcp.Enum("null", "ascending", "descending")),
	mcp.WithNumber("page", mcp.Description("page number, default 1"), mcp.DefaultNumber(1)),
	mcp.WithNumber("pageSize", mcp.Description("page size, default 50"), mcp.DefaultNumber(50)),
)

func ListContainersHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	req := map[string]interface{}{
		"page":     int(numArg(args, "page", 1)),
		"pageSize": int(numArg(args, "pageSize", 50)),
		"name":     strArg(args, "name", ""),
		"state":    strArg(args, "state", "all"),
		"orderBy":  strArg(args, "orderBy", "createdAt"),
		"order":    strArg(args, "order", "descending"),
	}
	res := map[string]interface{}{}
	return utils.NewPanelClient("POST", "/containers/search", utils.WithPayload(req)).Request(&res)
}
