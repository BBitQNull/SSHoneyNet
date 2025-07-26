package ls

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/BBitQNull/SSHoneyNet/core/dispatcher"
	"github.com/BBitQNull/SSHoneyNet/core/filesystem"
	fs_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client/fsclient"
	proc_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client/procclient"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/exescript"
)

type LsHandler struct {
	fsClient   fs_client.FSManageClient
	procClient proc_client.ProcManageClient
}

func NewLsHandler(procClient proc_client.ProcManageClient, fsClient fs_client.FSManageClient) *LsHandler {
	return &LsHandler{
		fsClient:   fsClient,
		procClient: procClient,
	}
}

func FormatLsOutput(files []filesystem.FileNodeInfo) string {
	if len(files) == 0 || files == nil {
		return ""
	}
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	for i := range files {
		fmt.Fprintf(w, "%s\t", files[i].Name)
		if (i+1)%4 == 0 {
			fmt.Fprint(w, "\n")
		}
	}
	if len(files)%4 != 0 {
		fmt.Fprint(w, "\n")
	}
	w.Flush()
	return buf.String()
}
func FormatLsLongOutput(files []filesystem.FileNodeInfo) string {
	if len(files) == 0 || files == nil {
		return ""
	}
	var buf bytes.Buffer
	// minwidth=0, tabwidth=0, padding=2, padchar=' ', flags=0（基本默认）
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	for _, f := range files {
		// 模拟权限字符串
		perm := "-"
		if f.IsDir {
			perm = "d"
		}
		perm += "rw-r--r--"

		// 固定用户/组为 root
		user := "root"
		group := "root"

		// 格式化时间
		t := time.Unix(f.ModTime, 0)
		timeStr := t.Format("Jan _2 15:04") // e.g., Jul 26 19:25

		// 写入格式化输出，\t 用于对齐
		fmt.Fprintf(w, "%s\t1\t%s\t%s\t%d\t%s\t%s\n",
			perm, user, group, f.Size, timeStr, f.Name,
		)
	}

	w.Flush()
	return buf.String()
}

func (h *LsHandler) Execute(ctx context.Context, cmd exescript.ExecCommand, sessionID string) (dispatcher.CmdEcho, error) {
	if cmd.Name != "ls" {
		return dispatcher.CmdEcho{
			CmdResult: "zsh: command not found: " + cmd.Name,
		}, nil
	}
	args := cmd.Args
	if len(args) == 0 {
		args = []string{"/"}
	}
	useLongFormat := false
	if cmd.Flags != nil {
		if _, ok := cmd.Flags["l"]; ok {
			useLongFormat = true
		}
	}

	var output strings.Builder
	for _, path := range args {
		resp, err := h.fsClient.ListChildren(ctx, &fs_client.RawFSRequest{
			Path: path,
		})
		if err != nil {
			log.Printf("ls: cannot access '%s': No such file or directory\n", path)
			output.WriteString(fmt.Sprintf("ls: cannot access '%s': No such file or directory\n", path))
			continue
		}

		v, ok := resp.(*fs_client.RawFSResponse)
		if !ok || v == nil {
			output.WriteString(fmt.Sprintf("ls: internal error for path '%s'\n", path))
			continue
		}

		// 如果多个路径，显示类似 `path:` 前缀
		if len(args) > 1 {
			output.WriteString(fmt.Sprintf("%s:\n", path))
		}
		formatted := ""
		if useLongFormat {
			formatted = FormatLsLongOutput(v.Children)
		} else {
			formatted = FormatLsOutput(v.Children)
		}
		if formatted != "" {
			output.WriteString(formatted)
		}
	}
	return dispatcher.CmdEcho{
		CmdResult: output.String(),
	}, nil
}
