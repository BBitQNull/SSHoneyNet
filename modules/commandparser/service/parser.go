package parser_service

import (
	"context"
	"fmt"
	"log"

	"github.com/BBitQNull/SSHoneyNet/core/commandparser"
	"github.com/alecthomas/participle/v2"
)

type cmdParserService struct {
	parser *participle.Parser[commandparser.Script]
}

func NewCmdParserService() (commandparser.CmdParserService, error) {
	parser, err := participle.Build[commandparser.Script]()
	if err != nil {
		return nil, fmt.Errorf("failed to build parser: %w", err)
	}
	return &cmdParserService{parser: parser}, nil
}

func (c *cmdParserService) CommandParser(ctx context.Context, cmd string) (*commandparser.Script, error) {
	ast, err := c.parser.ParseString("", cmd)
	if err != nil {
		log.Printf("failed to parse cmd: %v", err)
		return nil, err
	}
	return ast, nil
}
