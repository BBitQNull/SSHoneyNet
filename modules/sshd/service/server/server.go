package server

import (
	"flag"
	"log"
	"os"

	"github.com/BBitQNull/SSHoneyNet/modules/sshd/service/handler"
	"github.com/BBitQNull/SSHoneyNet/pkg/model"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/auth"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

var serverAddr = flag.String("addr", "127.0.0.1:2222", "server's address")

func loadPrivateKey(path string, passphrase []byte) (gossh.Signer, error) {
	dir, _ := os.Getwd()
	keyBytes, err := os.ReadFile(dir + path)
	if err != nil {
		return nil, err
	}
	return gossh.ParsePrivateKeyWithPassphrase(keyBytes, passphrase)
}

// 闭包 依赖于authSvc的具体实现
func StartServer(authSvc auth.AuthService) {
	signer, err := loadPrivateKey("/pkg/key/host_key_rsa", []byte("123456"))
	if err != nil {
		log.Fatal("failed to loadkeyfile: ", err)
	}
	s := &ssh.Server{
		Addr:    *serverAddr,
		Handler: handler.SessionHandler(model.NewEchoRegistry()),
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			username := ctx.User()
			return authSvc.PasswordValidator(username, password)
		},
		HostSigners: []ssh.Signer{signer},
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal("failed to start SSH server: ", err)
	}
}
