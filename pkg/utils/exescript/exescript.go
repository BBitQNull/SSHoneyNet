package exescript

import (
	"strings"

	"github.com/BBitQNull/SSHoneyNet/core/commandparser"
)

type ExecScript struct {
	Lines []ExecCommandLine
}

type ExecCommandLine struct {
	Pipeline []ExecCommand
	Redir    *ExecRedirection
}

type ExecCommand struct {
	Name  string            // 命令名，如 "ls"
	Flags map[string]string // 参数，如 {"-l": "", "--color": "auto"}
	Args  []string          // 参数，如 ["dir1", "dir2"]
}

type ExecRedirection struct {
	Type string // ">", ">>", "<", "<<"
	File string // 重定向目标文件
}

func ConvertScript(ast *commandparser.Script) *ExecScript {
	exec := &ExecScript{}
	for _, line := range ast.Lines {
		execLine := ExecCommandLine{}
		for _, cmd := range line.Pipeline {
			execCmd := ExecCommand{
				Name:  cmd.Name,
				Flags: make(map[string]string),
				Args:  []string{},
			}
			for _, f := range cmd.Flags {
				if f.Value != nil {
					execCmd.Flags[f.Name] = *f.Value
				} else {
					execCmd.Flags[f.Name] = ""
				}
			}
			for _, arg := range cmd.Args {
				if arg.Value != nil {
					execCmd.Args = append(execCmd.Args, *arg.Value)
				}
			}
			execLine.Pipeline = append(execLine.Pipeline, execCmd)
		}
		if line.Redir != nil {
			execLine.Redir = &ExecRedirection{
				Type: line.Redir.File[:1], // 提取 >、< 等
				File: strings.TrimPrefix(line.Redir.File, line.Redir.File[:1]),
			}
		}
		exec.Lines = append(exec.Lines, execLine)
	}
	return exec
}
