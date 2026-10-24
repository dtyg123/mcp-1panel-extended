package extended

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/dtyg123/mcp-1panel-extended/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	GetContainerLogs = "get_container_logs"
)

// maxLogLines caps how many log lines are returned to keep tool output bounded.
const maxLogLines = 5000

var GetContainerLogsTool = mcp.NewTool(
	GetContainerLogs,
	mcp.WithDescription("get logs of a docker container. Returns the last N log lines as plain text. Use list_containers to find container names."),
	mcp.WithString("container", mcp.Description("container name, e.g. php, openresty, 1panel-core"), mcp.Required()),
	mcp.WithNumber("tail", mcp.Description("number of last log lines to return, default 200"), mcp.DefaultNumber(200)),
	mcp.WithBoolean("follow", mcp.Description("when true, keep streaming new logs for up to ~25 seconds (the panel client enforces a 30s HTTP timeout)"), mcp.DefaultBool(false)),
	mcp.WithBoolean("timestamp", mcp.Description("prepend a timestamp to each line"), mcp.DefaultBool(false)),
	mcp.WithString("since", mcp.Description("only show logs newer than this, e.g. 10m, 1h, or an RFC3339 timestamp. optional"), mcp.DefaultString("")),
)

func GetContainerLogsHandle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments

	container, err := requiredStrArg(args, "container")
	if err != nil {
		return mcp.NewToolResultText(err.Error()), nil
	}

	tail := int(numArg(args, "tail", 200))
	if tail <= 0 {
		tail = 200
	}
	follow := boolArg(args, "follow", false)
	useTimestamp := boolArg(args, "timestamp", false)
	since := strArg(args, "since", "")

	query := map[string]interface{}{
		"container": container,
		"tail":      tail,
		"follow":    follow,
		"timestamp": useTimestamp,
	}
	if since != "" {
		query["since"] = since
	}

	client := utils.NewPanelClient("GET", "/containers/search/log",
		utils.WithQuery(query),
		// the endpoint only streams when the client asks for SSE
		utils.WithHeaders(map[string]string{"Accept": "text/event-stream"}),
	)

	if _, err := client.Do(); err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("failed to fetch logs for %s: %v", container, err)), nil
	}
	defer client.Response.Body.Close()

	// The panel streams Server-Sent Events: each log line arrives as
	// "data: <line>\n\n". Strip the prefix and join into plain text.
	scanner := bufio.NewScanner(client.Response.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var lines []string
	truncated := false
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		line = strings.TrimPrefix(line, "data:")
		line = strings.TrimPrefix(line, " ")
		if line == "" {
			continue
		}
		if len(lines) >= maxLogLines {
			truncated = true
			break
		}
		lines = append(lines, line)
	}

	readErr := scanner.Err()
	if readErr != nil && len(lines) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("failed to read logs for %s: %v", container, readErr)), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "== %s: last %d lines ==\n", container, len(lines))
	sb.WriteString(strings.Join(lines, "\n"))
	if truncated {
		fmt.Fprintf(&sb, "\n[output truncated at %d lines]", maxLogLines)
	}
	if follow && readErr != nil {
		// expected when the 30s HTTP client timeout closes a following stream
		fmt.Fprintf(&sb, "\n[stream ended: %v]", readErr)
	}
	return mcp.NewToolResultText(sb.String()), nil
}
