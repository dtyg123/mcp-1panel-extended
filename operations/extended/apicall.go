package extended

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dtyg123/mcp-1panel-extended/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	CallPanelAPI = "call_panel_api"
)

var CallPanelAPITool = mcp.NewTool(
	CallPanelAPI,
	mcp.WithDescription("call ANY 1Panel API endpoint directly. Paths are relative to /api/v2 (e.g. /containers/search, /runtimes/search, /hosts/firewall/rules). Use this for endpoints not covered by other tools. Read-only methods (GET/POST search-style queries) are safest; be careful with write operations."),
	mcp.WithString("method", mcp.Description("HTTP method"), mcp.DefaultString("GET"), mcp.Enum("GET", "POST", "PUT", "DELETE"), mcp.Required()),
	mcp.WithString("path", mcp.Description("API path relative to /api/v2, e.g. /containers/search or /runtimes/search. May include a ?query string."), mcp.Required()),
	mcp.WithString("body", mcp.Description("raw JSON request body for POST/PUT, e.g. {\"page\":1,\"pageSize\":10}. Leave empty for GET."), mcp.DefaultString("")),
	mcp.WithString("query", mcp.Description("extra query string appended to the path, e.g. page=1&pageSize=10 (alternative to embedding it in path)"), mcp.DefaultString("")),
)

func CallPanelAPIHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments

	method := strings.ToUpper(strArg(args, "method", "GET"))
	path, err := requiredStrArg(args, "path")
	if err != nil {
		return mcp.NewToolResultText(err.Error()), nil
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if q := strArg(args, "query", ""); q != "" {
		if strings.Contains(path, "?") {
			path += "&" + q
		} else {
			path += "?" + q
		}
	}

	opts := []utils.Option{}
	if body := strArg(args, "body", ""); body != "" && method != "GET" && method != "DELETE" {
		var payload interface{}
		if err := json.Unmarshal([]byte(body), &payload); err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("invalid body JSON: %v", err)), nil
		}
		opts = append(opts, utils.WithPayload(payload))
	}

	res := map[string]interface{}{}
	return utils.NewPanelClient(method, path, opts...).Request(&res)
}
