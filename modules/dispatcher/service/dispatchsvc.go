package service

import (
	"context"
	"errors"
	"sync"

	"github.com/BBitQNull/SSHoneyNet/core/commandparser"
	"github.com/BBitQNull/SSHoneyNet/core/dispatcher"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/exescript"
)

type CmdDispatcherServer struct{}

type CmdHandler interface {
	Execute(ast exescript.ExecCommand) (dispatcher.CmdEcho, error)
}

var (
	commandMap = make(map[string]CmdHandler)
	mu         sync.RWMutex
)

func RegisterCmd(name string, handler CmdHandler) {
	mu.Lock()
	defer mu.Unlock()
	commandMap[name] = handler
}

func ExecuteScript(ir exescript.ExecScript) (dispatcher.CmdEcho, error) {
	mu.Lock()
	defer mu.Unlock()
	do := false
	for _, item := range ir.Lines {
		for _, comment := range item.Pipeline {
			if !do {
				handler, ok := commandMap[comment.Name]
				if !ok {
					return dispatcher.CmdEcho{
						CmdResult: "",
						ErrCode:   1,
						ErrMsg:    "command not found: " + comment.Name,
					}, errors.New("command not found:")
				}
				result, err := handler.Execute(exescript.ExecCommand{
					Name:  comment.Name,
					Flags: comment.Flags,
					Args:  comment.Args,
				})
				if err != nil {
					// 临时
					return dispatcher.CmdEcho{
						CmdResult: "",
						ErrCode:   2,
						ErrMsg:    err.Error(),
					}, err
				}
				return result, nil
			}
		}
	}
	// 占位
	return dispatcher.CmdEcho{}, nil
}
func (s *CmdDispatcherServer) CmdDispatcher(ctx context.Context, astReq commandparser.Script) (dispatcher.CmdEcho, error) {
	// 占位
	return dispatcher.CmdEcho{}, nil
}
