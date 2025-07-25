package dispatch_service

import (
	"context"
	"log"
	"sync"

	"github.com/BBitQNull/SSHoneyNet/core/clientset"
	"github.com/BBitQNull/SSHoneyNet/core/commandparser"
	"github.com/BBitQNull/SSHoneyNet/core/dispatcher"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/exescript"
	"google.golang.org/grpc/metadata"
)

type CmdDispatcherServer struct {
	clients    *clientset.ClientSet
	commandMap map[string]CmdHandler
	mu         sync.RWMutex
}

type CmdHandler interface {
	Execute(ctx context.Context, ast exescript.ExecCommand, sessionID string) (dispatcher.CmdEcho, error)
}

func NewDispatcherServer(clients *clientset.ClientSet) *CmdDispatcherServer {
	return &CmdDispatcherServer{
		clients:    clients,
		commandMap: make(map[string]CmdHandler),
	}
}

var (
	commandMap = make(map[string]CmdHandler)
	mu         sync.RWMutex
)

func (s *CmdDispatcherServer) RegisterCmd(name string, handler CmdHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.commandMap[name] = handler
}

func (s *CmdDispatcherServer) ExecuteScript(ctx context.Context, ir exescript.ExecScript, sessionID string) (dispatcher.CmdEcho, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	do := false
	for _, item := range ir.Lines {
		for _, comment := range item.Pipeline {
			if !do {
				handler, ok := s.commandMap[comment.Name]
				if !ok {
					return dispatcher.CmdEcho{
						CmdResult: "",
						ErrCode:   1,
						ErrMsg:    "zsh: command not found: " + comment.Name,
					}, nil
				}
				result, err := handler.Execute(ctx, exescript.ExecCommand{
					Name:  comment.Name,
					Flags: comment.Flags,
					Args:  comment.Args,
				}, sessionID)
				md, ok := metadata.FromIncomingContext(ctx)
				log.Printf("CmdHandler.Execute metadata: %+v, ok=%v", md, ok)
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

func (s *CmdDispatcherServer) CmdDispatcher(ctx context.Context, astReq commandparser.Script, sessionID string) (dispatcher.CmdEcho, error) {
	ir := exescript.ConvertScript(&astReq)
	return s.ExecuteScript(ctx, *ir, sessionID)
}
