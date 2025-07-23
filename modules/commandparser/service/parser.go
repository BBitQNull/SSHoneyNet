package parser_service

import (
	"context"
	"log"

	"github.com/BBitQNull/SSHoneyNet/core/commandparser"
	"github.com/alecthomas/participle/v2"
)

type cmdParserService struct{}

func NewCmdParserService() commandparser.CmdParserService {
	return &cmdParserService{}
}

func (c *cmdParserService) CommandParser(ctx context.Context, cmd string) (*commandparser.Script, error) {
	parser, err := participle.Build[commandparser.Script]()
	if err != nil {
		log.Fatal("failed to build parser:", err)
		return nil, err
	}
	ast, err := parser.ParseString("", cmd)
	if err != nil {
		log.Fatal("failed to parser cmd: ", err)
		return nil, err
	}
	return ast, nil
}
