package server

import (
	"flag"
	"log"

	"github.com/BBitQNull/SSHoneyNet/modules/sshd/service/handler"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/auth"
	"github.com/gliderlabs/ssh"
)

var serverAddr = flag.String("addr", "127.0.0.1:2222", "server's address")

// 闭包 依赖于authSvc的具体实现
func StartServer(authSvc auth.AuthService) {
	s := &ssh.Server{
		Addr:    *serverAddr,
		Handler: handler.SessionHandler,
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			username := ctx.User()
			return authSvc.PasswordValidator(username, password)
		},
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal("failed to start SSH server: ", err)
	}
}
